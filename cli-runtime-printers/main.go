package main

import (
	"fmt"
	"os"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/kubernetes/scheme"
)

func main() {
	obj := &corev1.ConfigMap{
		Data: map[string]string{"foo": "bar"},
	}
	obj.Name = "my-cm"

	// YAML
	fmt.Println("# YAML ConfigMap representation")
	printr := printers.NewTypeSetter(scheme.Scheme).ToPrinter(&printers.YAMLPrinter{})
	if err := printr.PrintObj(obj, os.Stdout); err != nil {
		panic(err.Error())
	}

	fmt.Println()

	// JSON
	fmt.Println("# JSON ConfigMap representation")
	printr = printers.NewTypeSetter(scheme.Scheme).ToPrinter(&printers.JSONPrinter{})
	if err := printr.PrintObj(obj, os.Stdout); err != nil {
		panic(err.Error())
	}

	fmt.Println()

	// Table (human-readable)
	fmt.Println("# Table ConfigMap representation")
	printr = printers.NewTypeSetter(scheme.Scheme).ToPrinter(printers.NewTablePrinter(printers.PrintOptions{}))
	if err := printr.PrintObj(obj, os.Stdout); err != nil {
		panic(err.Error())
	}

	fmt.Println()

	// JSONPath
	fmt.Println("# ConfigMap.data.foo")
	printr, err := printers.NewJSONPathPrinter("{.data.foo}")
	if err != nil {
		panic(err.Error())
	}

	printr = printers.NewTypeSetter(scheme.Scheme).ToPrinter(printr)
	if err := printr.PrintObj(obj, os.Stdout); err != nil {
		panic(err.Error())
	}

	fmt.Println()

	// Name-only
	fmt.Println("# <kind>/<name>")
	printr = printers.NewTypeSetter(scheme.Scheme).ToPrinter(&printers.NamePrinter{})
	if err := printr.PrintObj(obj, os.Stdout); err != nil {
		panic(err.Error())
	}

	fmt.Println()
}
