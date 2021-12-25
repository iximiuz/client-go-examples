package main

import (
	"fmt"

	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{
			ExplicitPath: "/home/vagrant/.kube/config",
		},
		&clientcmd.ConfigOverrides{},
	).RawConfig()
	if err != nil {
		panic(err.Error())
	}

	for name := range config.Contexts {
		fmt.Printf("Found context %s\n", name)
	}
}
