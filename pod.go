package main

// NODEGROUP      NODE                            NAMESPACE        NAME                                                        READY   STATUS             RESTARTS   AGE     IP
// staging-spot   ip-10-100-29-185.ec2.internal   airflow          airflow-db-migrations-696f99f8b5-j72q5                      1/1     Running            0          23d     10.100.29.51
import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

type Pod struct {
	Nodegroup string
	Node      string
	Namespace string
	Name      string
	Ready     string
	Status    string
	Restarts  string
	Age       string
	IP        string
	Mem       int64
	Cpu       int64
}

func PodIsReady(pod corev1.Pod) string {
	running := 0
	for _, item := range pod.Status.ContainerStatuses {
		if item.Ready {
			running = running + 1
		}
	}
	return fmt.Sprintf("%d/%d", running, len(pod.Status.ContainerStatuses))
}

func GetContainerRestarts(pod corev1.Pod) string {
	restarts := int32(0)
	for _, item := range pod.Status.ContainerStatuses {
		restarts = restarts + item.RestartCount
	}
	return fmt.Sprintf("%d", restarts)
}
func GetPodAge(pod corev1.Pod) string {
	return "0" //fmt.Sprintf("%.fd", time.Since(pod.Status.StartTime.Time).Hours()/24)
}

func NewPod(pod corev1.Pod, nodegroup string) Pod {
	return Pod{
		Nodegroup: nodegroup,
		Node:      pod.Spec.NodeName,
		Namespace: pod.Namespace,
		Name:      pod.Name,
		Ready:     PodIsReady(pod),
		Status:    fmt.Sprintf("%v", pod.Status.Phase),
		Restarts:  GetContainerRestarts(pod),
		Age:       GetPodAge(pod),
		IP:        pod.Status.PodIP,
		Mem:       0,
		Cpu:       0,
	}
}

func (pod *Pod) String() []string {
	return []string{
		pod.Nodegroup,
		pod.Node,
		pod.Namespace,
		pod.Name,
		pod.Ready,
		pod.Status,
		pod.Restarts,
		pod.Age,
		pod.IP,
	}
}

func (pod *Pod) UsageString() []string {
	return []string{
		pod.Nodegroup,
		pod.Node,
		pod.Namespace,
		pod.Name,
		fmt.Sprintf("%vm", pod.Cpu),
		fmt.Sprintf("%vMi", pod.Mem),
	}
}

func PodHeader() []string {
	return []string{
		"Nodegroup",
		"Node",
		"Namespace",
		"Name",
		"Ready",
		"Status",
		"Restarts",
		"Age",
		"IP",
		"Cpu",
		"Mem",
	}
}

func PodUsageHeader() []string {
	return []string{
		"Nodegroup",
		"Node",
		"Namespace",
		"Name",
		"Cpu(Cores)",
		"Mem(Bytes)",
	}
}
