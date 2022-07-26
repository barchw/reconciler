package data

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/avast/retry-go"
	iopv1alpha1 "istio.io/istio/operator/pkg/apis/istio/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

//go:generate mockery --name=Gatherer --outpkg=mocks --case=underscore
// Gatherer gathers data from the Kubernetes cluster.
type Gatherer interface {
	// GetAllPods from the cluster and return them as a v1.PodList.
	GetAllPods(kubeClient kubernetes.Interface, retryOpts []retry.Option) (podsList *v1.PodList, err error)

	// GetPodsWithDifferentImage than the passed expected image to filter them out from the pods list.
	GetPodsWithDifferentImage(inputPodsList v1.PodList, image ExpectedImage) (outputPodsList v1.PodList)

	GetIstioOperator(dynamicClient dynamic.Interface) (*iopv1alpha1.IstioOperator, bool, error)
}

// DefaultGatherer that gets pods from the Kubernetes cluster
type DefaultGatherer struct{}

// ExpectedImage to be verified by the proxy.
type ExpectedImage struct {
	Prefix  string
	Version string
}

// NewDefaultGatherer creates a new instance of DefaultGatherer.
func NewDefaultGatherer() *DefaultGatherer {
	return &DefaultGatherer{}
}

func (i *DefaultGatherer) GetAllPods(kubeClient kubernetes.Interface, retryOpts []retry.Option) (podsList *v1.PodList, err error) {
	err = retry.Do(func() error {
		podsList, err = kubeClient.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
		if err != nil {
			return err
		}

		return nil
	}, retryOpts...)

	if err != nil {
		return nil, err
	}

	return
}

func (i *DefaultGatherer) GetPodsWithDifferentImage(inputPodsList v1.PodList, image ExpectedImage) (outputPodsList v1.PodList) {
	inputPodsList.DeepCopyInto(&outputPodsList)
	outputPodsList.Items = []v1.Pod{}

	for _, pod := range inputPodsList.Items {
		if _, containsIstioSidecarAnnotation := pod.Annotations["sidecar.istio.io/status"]; !containsIstioSidecarAnnotation || !isPodReady(pod) {
			continue
		}

		istioSidecarNames := getIstioSidecarNamesFromAnnotations(pod.Annotations)

		for _, container := range pod.Spec.Containers {
			if !isIstioSidecar(istioSidecarNames, container.Name) {
				continue
			}
			containsPrefix := strings.Contains(container.Image, image.Prefix)
			hasSuffix := strings.HasSuffix(container.Image, image.Version)
			if !hasSuffix || !containsPrefix {
				outputPodsList.Items = append(outputPodsList.Items, *pod.DeepCopy())
			}
		}
	}

	return
}

func (d *DefaultGatherer) GetIstioOperator(dynamicClient dynamic.Interface) (*iopv1alpha1.IstioOperator, bool, error) {
	operators := []iopv1alpha1.IstioOperator{}
	list, err := dynamicClient.Resource(schema.GroupVersionResource{Group: "install.istio.io", Version: "v1alpha1", Resource: "istiooperators"}).List(context.Background(), "installed-state-default-operator", metav1.GetOptions{})
	if err != nil {
		return nil, false, err
	}
	if len(list.Items) == 0 {
		return nil, false, nil
	}
	if len(list.Items) > 1 {
		return nil, true, errors.New("more than one IstioOperator found")
	}
	toUnmarshal, err := list.MarshalJSON()
	if err != nil {
		return nil, false, err
	}
	json.Unmarshal(toUnmarshal, &operators)

	return &operators[0], true, nil
}

// getIstioSidecarNamesFromAnnotations gets all container names in pod annoted with podAnnotations that are Istio sidecars
func getIstioSidecarNamesFromAnnotations(podAnnotations map[string]string) []string {
	type istioStatusStruct struct {
		Containers []string `json:"containers"`
	}
	istioStatus := istioStatusStruct{}
	err := json.Unmarshal([]byte(podAnnotations["sidecar.istio.io/status"]), &istioStatus)
	if err != nil {
		return []string{}
	}
	return istioStatus.Containers
}

// isIstioSidecar checks whether the pod with name=containerName is a Istio sidecar in pod with Istio sidecars with names=istioSidecarNames
func isIstioSidecar(istioSidecarNames []string, containerName string) bool {
	for _, c := range istioSidecarNames {
		if c == containerName {
			return true
		}
	}
	return false
}

// isPodReady checks if the pod is Ready, returns true if the Pod is in the Running state and not Pending or Terminating.
func isPodReady(pod v1.Pod) bool {

	if pod.Status.Phase != v1.PodRunning {
		return false
	}
	for _, condition := range pod.Status.Conditions {
		if condition.Status != v1.ConditionTrue {
			return false
		}
	}

	return pod.ObjectMeta.DeletionTimestamp == nil
}

// RemoveAnnotatedPods removes pods with annotation annotationKey from in podList
func RemoveAnnotatedPods(in v1.PodList, annotationKey string) (out v1.PodList) {
	in.DeepCopyInto(&out)
	out.Items = []v1.Pod{}
	for i := 0; i < len(in.Items); i++ {
		if _, ok := in.Items[i].Annotations[annotationKey]; !ok {
			out.Items = append(out.Items, in.Items[i])
		}
	}
	return
}
