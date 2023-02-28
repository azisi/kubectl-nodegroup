package main

import (
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned/fake"
)

func TestGetPods(t *testing.T) {

	n1 := &corev1.Node{
		ObjectMeta: v1.ObjectMeta{
			Name: "test-node-name-1",
			CreationTimestamp: v1.Time{
				Time: time.Date(2023, 1, 2, 3, 30, 23, 0, time.UTC),
			},
			Labels: map[string]string{"eks.amazonaws.com/nodegroup": "testgroup"},
		},
		Status: corev1.NodeStatus{
			NodeInfo: corev1.NodeSystemInfo{
				KubeletVersion: "v1.24",
			},
		},
	}
	n2 := &corev1.Node{
		ObjectMeta: v1.ObjectMeta{
			Name: "test-node-name-2",
			CreationTimestamp: v1.Time{
				Time: time.Date(2023, 1, 2, 3, 30, 23, 0, time.UTC),
			},
			Labels: map[string]string{"eks.amazonaws.com/nodegroup": "testgroup2"},
		},
		Status: corev1.NodeStatus{
			NodeInfo: corev1.NodeSystemInfo{
				KubeletVersion: "v1.24",
			},
		},
	}

	objs := []runtime.Object{n1, n2}
	client := fake.NewSimpleClientset(objs...)
	metricsclientset := metricsv.NewSimpleClientset()

	cfg = Config{
		kubeconfig:      StringP(""),
		filterNamespace: StringP(""),
		filterNodegroup: StringP(""),
		showWide:        BoolP(false),
		showPods:        BoolP(false),
		showUsage:       BoolP(false),
	}
	c := NewClient(client.CoreV1(), metricsclientset.MetricsV1beta1())
	c.ParseK8sData()
	c.PrintResults()
}
func BoolP(b bool) *bool {
	return &b
}
func StringP(s string) *string {
	return &s
}
