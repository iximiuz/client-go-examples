package main

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	namespace = "default"
	label     = "simple-list-typed-" + rand.String(6)
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

	// Create one object before starting to watch.
	first := createConfigMap(client)

	// Start watching. Expected events:
	//  - ADDED the first config map (even though it was done before starting the watch)
	//  - ADDED the second config map
	//  - x2 DELETED
	watch, err := client.
		CoreV1().
		ConfigMaps(namespace).
		Watch(
			context.Background(),
			metav1.ListOptions{
				LabelSelector: "example==" + label,
			},
		)
	if err != nil {
		panic(err.Error())
	}
	go func() {
		for event := range watch.ResultChan() {
			fmt.Printf(
				"Watch Event: %s %s\n",
				event.Type, event.Object.GetObjectKind().GroupVersionKind().Kind,
			)
		}
	}()

	// Create another object while watching.
	second := createConfigMap(client)

	deleteConfigMap(client, first)
	deleteConfigMap(client, second)

	time.Sleep(2 * time.Second)
	watch.Stop()
	time.Sleep(1 * time.Second)
}

func createConfigMap(client kubernetes.Interface) *corev1.ConfigMap {
	cm := &corev1.ConfigMap{Data: map[string]string{"foo": "bar"}}
	cm.Namespace = namespace
	cm.GenerateName = "simple-list-typed-"
	cm.SetLabels(map[string]string{"example": label})

	cm, err := client.
		CoreV1().
		ConfigMaps(namespace).
		Create(
			context.Background(),
			cm,
			metav1.CreateOptions{},
		)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Created ConfigMap %s/%s\n", cm.GetNamespace(), cm.GetName())
	return cm
}

func deleteConfigMap(client kubernetes.Interface, cm *corev1.ConfigMap) {
	err := client.
		CoreV1().
		ConfigMaps(cm.GetNamespace()).
		Delete(
			context.Background(),
			cm.GetName(),
			metav1.DeleteOptions{},
		)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Deleted ConfigMap %s/%s\n", cm.GetNamespace(), cm.GetName())
}
