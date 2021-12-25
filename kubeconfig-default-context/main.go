package main

import (
	"fmt"

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", "/home/vagrant/.kube/config")
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
