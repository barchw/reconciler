package service

import (
	"context"
	"github.com/pkg/errors"
	"time"

	"github.com/kyma-incubator/reconciler/pkg/model"
	"go.uber.org/zap"
)

const (
	defaultOperationsWatchInterval = 30 * time.Second
	defaultOrphanOperationTimeout  = 10 * time.Minute
)

type BookkeeperConfig struct {
	OperationsWatchInterval time.Duration
	OrphanOperationTimeout  time.Duration
}

func (wc *BookkeeperConfig) validate() error {
	if wc.OperationsWatchInterval < 0 {
		return errors.New("operations watch interval cannot be < 0")
	}
	if wc.OperationsWatchInterval == 0 {
		wc.OperationsWatchInterval = defaultOperationsWatchInterval
	}
	if wc.OrphanOperationTimeout < 0 {
		return errors.New("orphan operation timeout cannot be < 0")
	}
	if wc.OrphanOperationTimeout == 0 {
		wc.OrphanOperationTimeout = defaultOrphanOperationTimeout
	}
	return nil
}

type bookkeeper struct {
	transition *ClusterStatusTransition
	config     *BookkeeperConfig
	logger     *zap.SugaredLogger
}

func newBookkeeper(transition *ClusterStatusTransition, config *BookkeeperConfig, logger *zap.SugaredLogger) *bookkeeper {
	if config == nil {
		config = &BookkeeperConfig{}
	}
	return &bookkeeper{
		transition: transition,
		config:     config,
		logger:     logger,
	}
}

func (bk *bookkeeper) Run(ctx context.Context) error {
	if err := bk.config.validate(); err != nil {
		return err
	}

	bk.logger.Infof("Starting bookkeeper: interval for updating reconciliation statuses and orphan operations "+
		"is %.1f secs / timeout for orphan operations is %.1f secs",
		bk.config.OperationsWatchInterval.Seconds(), bk.config.OrphanOperationTimeout.Seconds())

	ticker := time.NewTicker(bk.config.OperationsWatchInterval)
	for {
		select {
		case <-ticker.C:
			ops, err := bk.transition.ReconciliationRepository().GetReconcilingOperations()
			if err != nil {
				bk.logger.Errorf("Failed to retrieve operations of currently running reconciliations: %s", err)
				continue
			}

			//define reconciliation status by checking all running operations
			filterResults, err := bk.processReconciliations(ops)
			if err != nil {
				bk.logger.Errorf("Processing of reconciliations statuses failed: %s", err)
				continue
			}

			//finish reconciliations
			for _, filterResult := range filterResults {
				newClusterStatus := filterResult.GetResult()
				bk.logger.Debugf("Bookkeeper evaluted cluster status of reconciliation '%s' to '%s' "+
					"(ops-statuses were: done=%d/error=%d/other=%d)",
					filterResult.schedulingID, newClusterStatus,
					len(filterResult.done), len(filterResult.error), len(filterResult.other))

				if newClusterStatus == model.ClusterStatusReady || newClusterStatus == model.ClusterStatusError {
					bk.logger.Infof("Bookkeeper is updating reconciliation '%s': final status is '%s'",
						filterResult.schedulingID, filterResult.GetResult())

					if err := bk.transition.FinishReconciliation(filterResult.schedulingID, newClusterStatus); err != nil {
						bk.logger.Errorf("Bookkeeper failed to update status of reconciliation "+
							"(schedulingID:%s) to '%s': %s", filterResult.schedulingID, newClusterStatus, err)
					}
				}

				//reset orphaned operations
				for _, orphanOp := range filterResult.GetOrphans() {
					if orphanOp.State == model.OperationStateOrphan {
						//don't update orphan operations which are already marked as 'orphan'
						continue
					}
					bk.logger.Infof("Bookkeeper is marking operation '%s' as orphan "+
						"(last update happened %.2f minutes ago)", orphanOp, time.Since(orphanOp.Updated).Minutes())

					if err := bk.transition.ReconciliationRepository().UpdateOperationState(
						orphanOp.SchedulingID, orphanOp.CorrelationID, model.OperationStateOrphan); err != nil {
						bk.logger.Errorf("Failed to update status of orphan operation '%s': %s", orphanOp, err)
					}
				}
			}
		case <-ctx.Done():
			bk.logger.Info("Stopping bookkeeper because parent context got closed")
			ticker.Stop()
			return nil
		}
	}
}

func (bk *bookkeeper) processReconciliations(ops []*model.OperationEntity) (map[string]*ReconciliationResult, error) {
	start := time.Now()
	reconStatuses := make(map[string]*ReconciliationResult)
	for _, op := range ops {
		reconStatus, ok := reconStatuses[op.SchedulingID]
		if !ok {
			reconStatus = newReconciliationResult(op.SchedulingID, bk.config.OrphanOperationTimeout, bk.logger)
			reconStatuses[op.SchedulingID] = reconStatus
		}
		if err := reconStatus.AddOperation(op); err != nil {
			return nil, err
		}
	}

	bk.logger.Infof("Bookkeeper processed %d running reconciliations with %d operations in %.1f secs",
		len(reconStatuses), len(ops), time.Since(start).Seconds())

	return reconStatuses, nil
}