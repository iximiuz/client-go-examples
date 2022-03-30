package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err.Error())
	}

	kubeconfigGetter := func() (*api.Config, error) {
		kubeconfigYAML, err := ioutil.ReadFile(path.Join(home, ".kube/config"))
		if err != nil {
			return nil, err
		}
		return clientcmd.Load([]byte(kubeconfigYAML))
	}

	config, err := clientcmd.BuildConfigFromKubeconfigGetter("", kubeconfigGetter)
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
