package main

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"

	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
	"k8s.io/metrics/pkg/client/clientset/versioned/typed/metrics/v1beta1"
)

func main() {

	ParseCLI()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *cfg.kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// create the metrics clientset
	metricsclientset := metricsv.NewForConfigOrDie(config)

	c := NewClient(clientset.CoreV1(), metricsclientset.MetricsV1beta1())
	c.ParseK8sData()
	c.PrintResults()
}

type Client struct {
	k8s        v1.CoreV1Interface
	k8sMetrics v1beta1.MetricsV1beta1Interface
	nl         *NodeList
	pl         *PodList
}

func NewClient(k8s v1.CoreV1Interface, k8sMetrics v1beta1.MetricsV1beta1Interface) *Client {
	return &Client{
		k8s:        k8s,
		k8sMetrics: k8sMetrics,
		nl:         NewNodeList(),
		pl:         NewPodList(),
	}
}

func (c *Client) ParseK8sData() {
	// get all nodes
	nodes, err := c.k8s.Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	// parse node information
	for _, node := range nodes.Items {
		c.nl.Append(node)
	}
	// get all pods
	pods, err := c.k8s.Pods(*cfg.filterNamespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	// parse pod information
	for _, pod := range pods.Items {
		c.pl.Append(pod, c.nl.GetNodeGroup(pod.Spec.NodeName))
		c.nl.StatsPodInc(c.nl.GetNodeGroup(pod.Spec.NodeName))
	}
	// get all pod metrics
	podMetrics, err := c.k8sMetrics.PodMetricses(*cfg.filterNamespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	// parse pod metrics
	for _, podmetric := range podMetrics.Items {
		c.pl.AddMetric(podmetric)
		c.nl.AddMetric(podmetric, c.pl.GetNode(podmetric.Name))
	}
}

func (c *Client) PrintResults() {
	if !*cfg.showPods { // don't show pods
		if *cfg.showWide { // show nodes
			c.nl.Print(*cfg.filterNodegroup)
		} else { // else show nodegroups
			c.nl.PrintNodeGroup(*cfg.filterNodegroup)
		}
	} else {
		if !*cfg.showUsage { // show pods
			c.pl.Print()
		} else { // show pod usage
			c.pl.PrintUsage()
		}
	}
}
