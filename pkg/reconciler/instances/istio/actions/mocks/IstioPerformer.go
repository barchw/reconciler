// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	chart "github.com/kyma-incubator/reconciler/pkg/reconciler/chart"
	actions "github.com/kyma-incubator/reconciler/pkg/reconciler/instances/istio/actions"

	context "context"

	helpers "github.com/kyma-incubator/reconciler/pkg/reconciler/instances/istio/helpers"

	kubernetes "github.com/kyma-incubator/reconciler/pkg/reconciler/kubernetes"

	mock "github.com/stretchr/testify/mock"

	zap "go.uber.org/zap"
)

// IstioPerformer is an autogenerated mock type for the IstioPerformer type
type IstioPerformer struct {
	mock.Mock
}

// Install provides a mock function with given fields: kubeConfig, istioChart, version, logger
func (_m *IstioPerformer) Install(kubeConfig string, istioChart string, version helpers.HelperVersion, logger *zap.SugaredLogger) error {
	ret := _m.Called(kubeConfig, istioChart, version, logger)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, helpers.HelperVersion, *zap.SugaredLogger) error); ok {
		r0 = rf(kubeConfig, istioChart, version, logger)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PatchMutatingWebhook provides a mock function with given fields: ctx, kubeClient, logger
func (_m *IstioPerformer) PatchMutatingWebhook(ctx context.Context, kubeClient kubernetes.Client, logger *zap.SugaredLogger) error {
	ret := _m.Called(ctx, kubeClient, logger)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, kubernetes.Client, *zap.SugaredLogger) error); ok {
		r0 = rf(ctx, kubeClient, logger)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ResetProxy provides a mock function with given fields: _a0, kubeConfig, proxyImageVersion, logger
func (_m *IstioPerformer) ResetProxy(_a0 context.Context, kubeConfig string, proxyImageVersion helpers.HelperVersion, logger *zap.SugaredLogger) error {
	ret := _m.Called(_a0, kubeConfig, proxyImageVersion, logger)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, helpers.HelperVersion, *zap.SugaredLogger) error); ok {
		r0 = rf(_a0, kubeConfig, proxyImageVersion, logger)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Uninstall provides a mock function with given fields: kubeClientSet, version, logger
func (_m *IstioPerformer) Uninstall(kubeClientSet kubernetes.Client, version helpers.HelperVersion, logger *zap.SugaredLogger) error {
	ret := _m.Called(kubeClientSet, version, logger)

	var r0 error
	if rf, ok := ret.Get(0).(func(kubernetes.Client, helpers.HelperVersion, *zap.SugaredLogger) error); ok {
		r0 = rf(kubeClientSet, version, logger)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: kubeConfig, istioChart, targetVersion, logger
func (_m *IstioPerformer) Update(kubeConfig string, istioChart string, targetVersion helpers.HelperVersion, logger *zap.SugaredLogger) error {
	ret := _m.Called(kubeConfig, istioChart, targetVersion, logger)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, helpers.HelperVersion, *zap.SugaredLogger) error); ok {
		r0 = rf(kubeConfig, istioChart, targetVersion, logger)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Version provides a mock function with given fields: workspace, branchVersion, istioChart, kubeConfig, logger
func (_m *IstioPerformer) Version(workspace chart.Factory, branchVersion string, istioChart string, kubeConfig string, logger *zap.SugaredLogger) (actions.IstioStatus, error) {
	ret := _m.Called(workspace, branchVersion, istioChart, kubeConfig, logger)

	var r0 actions.IstioStatus
	if rf, ok := ret.Get(0).(func(chart.Factory, string, string, string, *zap.SugaredLogger) actions.IstioStatus); ok {
		r0 = rf(workspace, branchVersion, istioChart, kubeConfig, logger)
	} else {
		r0 = ret.Get(0).(actions.IstioStatus)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(chart.Factory, string, string, string, *zap.SugaredLogger) error); ok {
		r1 = rf(workspace, branchVersion, istioChart, kubeConfig, logger)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewIstioPerformer interface {
	mock.TestingT
	Cleanup(func())
}

// NewIstioPerformer creates a new instance of IstioPerformer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIstioPerformer(t mockConstructorTestingTNewIstioPerformer) *IstioPerformer {
	mock := &IstioPerformer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}