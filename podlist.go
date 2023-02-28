package main

import (
	"fmt"

	"github.com/azisi/tableformat"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

type PodList struct {
	Pods map[string]*Pod
}

func NewPodList() *PodList {
	return &PodList{
		Pods: make(map[string]*Pod),
	}
}

func (p *PodList) Append(pod corev1.Pod, nodegroup string) {
	tmp := NewPod(pod, nodegroup)

	if _, ok := p.Pods[tmp.Name]; ok {
		fmt.Printf("ERROR duplicate pod name %v\n", tmp.Name)
	} else {
		p.Pods[tmp.Name] = &tmp
	}
}

func (pl *PodList) GetNodeGroup(podName string) string {
	if pod, ok := pl.Pods[podName]; ok {
		return pod.Nodegroup
	} else {
		fmt.Printf("ERROR getting node group, cannot find pod name %v\n", podName)
		return "N/A"
	}
}

func (pl *PodList) GetNode(podName string) string {
	if pod, ok := pl.Pods[podName]; ok {
		return pod.Node
	} else {
		fmt.Printf("ERROR getting node name, cannot find pod name %v\n", podName)
		return "N/A"
	}
}

func (n *PodList) Print() {
	nodeTable := tableformat.NewTable([]string{})
	for _, pod := range n.Pods {
		nodeTable.Append(pod.String())
	}
	nodeTable.OrderBy([]int{0, 1, 2})
	nodeTable.SetHeader(PodHeader())
	nodeTable.Print()
}

func (n *PodList) PrintUsage() {
	nodeTable := tableformat.NewTable([]string{})
	for _, pod := range n.Pods {
		nodeTable.Append(pod.UsageString())
	}
	nodeTable.OrderBy([]int{0, 1, 2})
	nodeTable.SetHeader(PodUsageHeader())
	nodeTable.Print()
}

func (pl *PodList) AddMetric(podmetric v1beta1.PodMetrics) {
	if pod, ok := pl.Pods[podmetric.Name]; ok {
		for _, containermetric := range podmetric.Containers {
			pod.Cpu += containermetric.Usage.Cpu().MilliValue()
			pod.Mem += containermetric.Usage.Memory().Value() / 1024 / 1024
			// fmt.Printf("%v %v %v %vMi\n", podmetric.Name, containermetric.Name, containermetric.Usage.Cpu().MilliValue(), containermetric.Usage.Memory().Value()/1024/1024)
		}
	} else {
		fmt.Printf("ERROR adding podmetric, cannot find pod name %v\n", podmetric.Name)
	}
}
