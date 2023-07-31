package main

import (
	"fmt"
	"os"
	"path"

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: path.Join(home, ".kube/config")},
		&clientcmd.ConfigOverrides{CurrentContext: os.Args[1]},
	).ClientConfig()
	if err != nil {
		panic(err.Error())
	}

	client, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	ver, err := client.ServerVersion()
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(ver.String())
}
