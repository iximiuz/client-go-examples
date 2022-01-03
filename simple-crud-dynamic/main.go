package main

import (
	"context"
	"fmt"
	"os"
	"path"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	config, err := clientcmd.BuildConfigFromFlags("", path.Join(home, ".kube/config"))
	if err != nil {
		panic(err.Error())
	}

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	namespace := "default"

	res := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "configmaps"}

	desired := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]interface{}{
				"namespace":    namespace,
				"generateName": "simple-crud-dynamic-",
			},
			"spec": map[string]interface{}{
				"data": map[string]interface{}{
					"foo": "bar",
				},
			},
		},
	}

	// Create
	created, err := client.
		Resource(res).
		Namespace(namespace).
		Create(context.Background(), desired, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Created ConfigMap %s/%s\n", namespace, created.GetName())

	// Read
	read, err := client.
		Resource(res).
		Namespace(namespace).
		Get(
			context.Background(),
			created.GetName(),
			metav1.GetOptions{},
		)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Read ConfigMap %s/%s\n", namespace, read.GetName())

	// Update
	updated, err := client.
		Resource(res).
		Namespace(namespace).
		Update(
			context.Background(),
			read,
			metav1.UpdateOptions{},
		)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Updated ConfigMap %s/%s\n", namespace, updated.GetName())

	// Delete
	err = client.
		Resource(res).
		Namespace(namespace).
		Delete(
			context.Background(),
			created.GetName(),
			metav1.DeleteOptions{},
		)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Deleted ConfigMap %s/%s\n", namespace, created.GetName())
}
