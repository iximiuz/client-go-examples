package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// go run main.go --help
// go run main.go --kubeconfig /foo/bar (panic: stat /foo/bar: no such file or directory)
// go run main.go --context shared1
// go run main.go --context shared2

func main() {
	configFlags := genericclioptions.NewConfigFlags(true)

	cmd := &cobra.Command{
		Use: "kubectl (well, almost)",
		Run: func(cmd *cobra.Command, args []string) {
			// Interesting methods:
			// - ToRawKubeConfigLoader (returns clientcmd.ClientConfig)
			// - ToRESTConfig
			// - ToRESTMapper
			// - ToDiscoveryClient

			kubeconfig, err := configFlags.ToRawKubeConfigLoader().RawConfig()
			if err != nil {
				panic(err.Error())
			}
			for name := range kubeconfig.Contexts {
				fmt.Printf("Found context %s\n", name)
			}

			restconfig, err := configFlags.ToRESTConfig()
			if err != nil {
				panic(err.Error())
			}
			fmt.Println("Cluster host", restconfig.Host)

			client, err := configFlags.ToDiscoveryClient()
			if err != nil {
				panic(err.Error())
			}

			ver, err := client.ServerVersion()
			if err != nil {
				panic(err.Error())
			}

			fmt.Println(ver.String())
		},
	}
	configFlags.AddFlags(cmd.PersistentFlags())

	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
