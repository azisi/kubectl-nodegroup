package main

import (
	"fmt"

	"github.com/azisi/tableformat"
	v1 "k8s.io/api/core/v1"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

type NodeGroupStat struct {
	Nodes int
	Pods  int
}

type NodeList struct {
	Nodes map[string]*Node
	Stats map[string]*NodeGroupStat
}

// Create a new node list
func NewNodeList() *NodeList {
	return &NodeList{
		Nodes: make(map[string]*Node),
		Stats: make(map[string]*NodeGroupStat),
	}
}

// Append a Node in the NodeList
func (n *NodeList) Append(node v1.Node) {
	tmp := NewNode(node)

	if _, ok := n.Nodes[tmp.Name]; ok {
		fmt.Printf("ERROR duplicate node name %v\n", tmp.Name)
	} else {
		n.Nodes[tmp.Name] = &tmp
	}
	if _, ok := n.Stats[tmp.Nodegroup]; ok {
		n.Stats[tmp.Nodegroup].Nodes++
	} else {
		n.Stats[tmp.Nodegroup] = &NodeGroupStat{Nodes: 1, Pods: 0}
	}
}

// Find a nodegroup in the NodeList and increase the pod count by 1
func (n *NodeList) StatsPodInc(nodegroup string) {
	if _, ok := n.Stats[nodegroup]; ok {
		n.Stats[nodegroup].Pods++
	} else {
		fmt.Printf("ERROR nodegroup %v not found, cannot increase pods\n", nodegroup)
	}
}

// Find a node in the NodeList and return its nodegroup
func (n *NodeList) GetNodeGroup(nodeName string) string {
	if node, ok := n.Nodes[nodeName]; ok {
		return node.Nodegroup
	} else {
		fmt.Printf("ERROR cannot find node group name for node %v\n", nodeName)
		return "N/A"
	}
}

// Print the nodes in a table format in stdout
func (n *NodeList) Print(filterNodegroup string) {
	nodeTable := tableformat.NewTable([]string{})
	for _, node := range n.Nodes {
		if filterNodegroup == "" || node.Nodegroup == filterNodegroup {
			nodeTable.Append(node.String())
		}
	}
	nodeTable.OrderBy([]int{0, 1, 2})
	nodeTable.SetHeader(NodeHeader(false))
	nodeTable.Print()
}

// Print the nodegroups in a table format in stdout
func (n *NodeList) PrintNodeGroup(filterNodegroup string) {
	nodeTable := tableformat.NewTable([]string{})
	nodeTable.OrderBy([]int{0, 1})
	nodeTable.SetHeader([]string{"Nodegroup", "Nodes", "Pods", "Pods/Node"})
	for name, stat := range n.Stats {
		if filterNodegroup == "" || name == filterNodegroup {
			nodeTable.Append([]string{name, fmt.Sprintf("%d", stat.Nodes), fmt.Sprintf("%d", stat.Pods), fmt.Sprintf("%d", stat.Pods/stat.Nodes)})
		}
	}
	nodeTable.Print()
}

// Find a node in the nodelist and add the pod cpu/mem metrics to the node metrics
func (nl *NodeList) AddMetric(podmetric v1beta1.PodMetrics, nodename string) {
	if node, ok := nl.Nodes[nodename]; ok {
		for _, containermetric := range podmetric.Containers {
			node.CpuUsage += containermetric.Usage.Cpu().MilliValue()
			node.MemUsage += containermetric.Usage.Memory().Value() / 1024 / 1024
			// fmt.Printf("%v %v %v %vMi\n", podmetric.Name, containermetric.Name, containermetric.Usage.Cpu().MilliValue(), containermetric.Usage.Memory().Value()/1024/1024)
		}
	} else {
		fmt.Printf("ERROR adding podmetric, cannot find node name %v\n", podmetric.Name)
	}
}
