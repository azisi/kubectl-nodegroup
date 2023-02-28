package main

import (
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
)

type Node struct {
	Nodegroup string
	Name      string
	Status    string
	Age       string
	Version   string
	CpuUsage  int64
	MemUsage  int64
}

func GetNodeIsReady(node corev1.Node) string {
	for _, item := range node.Status.Conditions {
		if item.Type == corev1.NodeReady && item.Status == corev1.ConditionTrue {
			return "True"
		}
	}
	return "False"
}

func GetNodeAge(node corev1.Node) string {
	return fmt.Sprintf("%.fd", time.Since(node.ObjectMeta.CreationTimestamp.Time).Hours()/24)
}

func GetNodeGroup(node corev1.Node) string {
	for key, value := range node.Labels {
		if key == "eks.amazonaws.com/nodegroup" {
			return value
		}
	}
	return "N/A"
}

func NewNode(node corev1.Node) Node {
	return Node{
		Nodegroup: GetNodeGroup(node),
		Name:      node.Name,
		Status:    GetNodeIsReady(node),
		Age:       GetNodeAge(node),
		Version:   node.Status.NodeInfo.KubeletVersion,
		CpuUsage:  0,
		MemUsage:  0,
	}
}

func (node *Node) String() []string {
	return []string{
		node.Nodegroup,
		node.Name,
		node.Status,
		node.Age,
		node.Version,
		fmt.Sprintf("%vm", node.CpuUsage),
		fmt.Sprintf("%vMiB", node.MemUsage),
	}
}
func NodeHeader(printUsage bool) []string {
	return []string{
		"Nodegroup",
		"Name",
		"Status",
		"Age",
		"Version",
		"Cpu(cores)",
		"Memory(bytes)",
	}
}
