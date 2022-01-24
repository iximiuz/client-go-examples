package main

import (
	"context"
	"fmt"
	"os"
	"path"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	// ERR_NOT_FOUND
	_, err = client.
		CoreV1().
		ConfigMaps(namespace).
		Get(
			context.Background(),
			"this_name_definitely_does_not_exist",
			metav1.GetOptions{},
		)
	if err == nil {
		panic("ERR_NOT_FOUND expected")
	}
	if !errors.IsNotFound(err) {
		panic(err.Error())
	}

	desired := corev1.ConfigMap{Data: map[string]string{"foo": "bar"}}
	desired.Namespace = namespace
	desired.GenerateName = "crud-typed-simple-"

	// Create
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

	// ERR_ALREADY_EXISTS
	duplicate := corev1.ConfigMap{}
	desired.Namespace = namespace
	duplicate.Name = created.Name
	_, err = client.
		CoreV1().
		ConfigMaps(namespace).
		Create(
			context.Background(),
			&duplicate,
			metav1.CreateOptions{},
		)
	if err == nil {
		panic("ERR_ALREADY_EXISTS expected")
	}
	if !errors.IsAlreadyExists(err) {
		panic(err.Error())
	}

	// Update
	created.Data["qux"] = "abc"
	updated, err := client.
		CoreV1().
		ConfigMaps(namespace).
		Update(
			context.Background(),
			created,
			metav1.UpdateOptions{},
		)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Updated ConfigMap %s/%s\n", namespace, updated.GetName())

	// ERR_CONFLICT
	created.Data["baz"] = "def"
	_, err = client.
		CoreV1().
		ConfigMaps(namespace).
		Update(
			context.Background(),
			created,
			metav1.UpdateOptions{},
		)
	if err == nil {
		panic("ERR_CONFLICT expected")
	}
	if !errors.IsConflict(err) {
		panic(err.Error())
	}
}
