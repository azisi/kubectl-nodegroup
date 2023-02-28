package main

import (
	"path/filepath"

	flag "github.com/spf13/pflag"
	"k8s.io/client-go/util/homedir"
)

var (
	cfg Config
)

type Config struct {
	kubeconfig      *string
	filterNamespace *string
	filterNodegroup *string
	showWide        *bool
	showPods        *bool
	showUsage       *bool
}

func ParseCLI() {
	cfg := &Config{}
	if home := homedir.HomeDir(); home != "" {
		cfg.kubeconfig = flag.StringP("kubeconfig", "k", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		cfg.kubeconfig = flag.StringP("kubeconfig", "k", "", "absolute path to the kubeconfig file")
	}
	cfg.filterNamespace = flag.StringP("namespace", "n", "", "Filter by namespace")
	cfg.filterNodegroup = flag.StringP("group", "g", "", "Filter by nodegroup")
	cfg.showWide = flag.BoolP("wide", "w", false, "Show more detailed information")
	cfg.showPods = flag.BoolP("pods", "p", false, "Show information regarding pods instead of nodes")
	cfg.showUsage = flag.BoolP("usage", "u", false, "Show information regarding cpu/mem usage instead of the default view")

	// flagHelp := flag.BoolP("help", "h", false, "Print this help message")

	flag.Parse()

	// if *flagHelp {
	// 	fmt.Fprintf(os.Stderr, "K8s deployment status\n\n")
	// 	fmt.Fprintf(os.Stderr, "Compare deployment versions across kubernetes clusters with latest commits and tags in Gitlab\n")
	// 	fmt.Fprintf(os.Stderr, "\n")
	// 	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	// 	flag.PrintDefaults()
	// 	fmt.Fprintf(os.Stderr, "\n")
	// 	fmt.Fprintf(os.Stderr, "The following precedence order is used for vars: cli -> conf -> defaults\n")
	// 	fmt.Fprintf(os.Stderr, "Config file name: %v.%v\n", cfgName, cfgType)
	// 	fmt.Fprintf(os.Stderr, "Config file path order: local directory first, then $HOME\n")
	// 	fmt.Fprintf(os.Stderr, "\n")
	// 	os.Exit(0)
	// }
}
