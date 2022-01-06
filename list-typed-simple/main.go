package main

import (
	"context"
	"fmt"
	"os"
	"path"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/client-go/kubernetes"
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

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	namespace := "default"
	label := "list-typed-simple-" + rand.String(6)

	desired := corev1.ConfigMap{Data: map[string]string{"foo": "bar"}}
	desired.Namespace = namespace
	desired.GenerateName = "list-typed-simple-"
	desired.SetLabels(map[string]string{"example": label})

	// Create a bunch of objects first.
	for i := 0; i < 10; i++ {
		created, err := client.
			CoreV1().
			ConfigMaps(namespace).
			Create(
				context.Background(),
				&desired,
				metav1.CreateOptions{},
			)
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("Created ConfigMap %s/%s\n", namespace, created.GetName())
	}

	// List - filter by the `example` label.
	list, err := client.
		CoreV1().
		ConfigMaps(namespace).
		List(
			context.Background(),
			metav1.ListOptions{
				LabelSelector: "example==" + label,
			},
		)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Found %d ConfigMap objects\n", len(list.Items))
}
