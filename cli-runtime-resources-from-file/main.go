package main

import (
	"os"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/kubernetes/scheme"
)

// go run main.go resources.yaml

func main() {
	configFlags := genericclioptions.NewConfigFlags(true)

	cmd := &cobra.Command{
		Use:  "kubectl (well, almost)",
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			builder := resource.NewBuilder(configFlags)

			namespace := ""
			if configFlags.Namespace != nil {
				namespace = *configFlags.Namespace
			}
			enforceNamespace := namespace != ""

			printr := printers.NewTypeSetter(scheme.Scheme).ToPrinter(&printers.YAMLPrinter{})

			err := builder.
				WithScheme(scheme.Scheme, scheme.Scheme.PrioritizedVersionsAllGroups()...).
				NamespaceParam(namespace).
				DefaultNamespace().
				FilenameParam(enforceNamespace, &resource.FilenameOptions{Filenames: args}).
				Do().
				Visit(func(info *resource.Info, err error) error {
					if err != nil {
						return err
					}

					return printr.PrintObj(info.Object, os.Stdout)
				})
			if err != nil {
				panic(err.Error())
			}

		},
	}
	configFlags.AddFlags(cmd.PersistentFlags())

	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
