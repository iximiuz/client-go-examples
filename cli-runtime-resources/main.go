package main

import (
	"os"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/kubernetes/scheme"
)

// go run main.go --help
// go run main.go po
// go run main.go pod
// go run main.go pods
// go run main.go services,deployments
// go run main.go --namespace=default service/kubernetes

func main() {
	configFlags := genericclioptions.NewConfigFlags(true)

	cmd := &cobra.Command{
		Use: "kubectl (well, almost)",
		Run: func(cmd *cobra.Command, args []string) {
			builder := resource.NewBuilder(configFlags)

			namespace := ""
			if configFlags.Namespace != nil {
				namespace = *configFlags.Namespace
			}

			obj, err := builder.
				WithScheme(scheme.Scheme, scheme.Scheme.PrioritizedVersionsAllGroups()...).
				NamespaceParam(namespace).
				DefaultNamespace().
				ResourceTypeOrNameArgs(true, args[0]).
				Do().
				Object()
			if err != nil {
				panic(err.Error())
			}

			printr := printers.NewTypeSetter(scheme.Scheme).ToPrinter(&printers.YAMLPrinter{})
			if err := printr.PrintObj(obj, os.Stdout); err != nil {
				panic(err.Error())
			}
		},
	}
	configFlags.AddFlags(cmd.PersistentFlags())

	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
