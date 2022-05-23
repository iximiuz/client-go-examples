package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"
)

var (
	namespace = "default"
	name      = "foobar"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	cfg, err := clientcmd.BuildConfigFromFlags("", path.Join(home, ".kube/config"))
	if err != nil {
		panic(err.Error())
	}

	client := kubernetes.NewForConfigOrDie(cfg)
	desired := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: map[string]string{"foo": "bar"},
	}

	defer deleteConfigMap(client, name, namespace)
	_, err = client.
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

	retry.RetryOnConflict(retry.DefaultRetry, func() error {
		one, err := getConfigMap(client, name, namespace)
		if err != nil {
			return err
		}
		// get point for configmap twice and only update one
		// then try to update two it will fail
		two, err := getConfigMap(client, name, namespace)
		if err != nil {
			return err
		}

		// generate random int so we're always different
		one.Data = map[string]string{
			"changed": fmt.Sprint(rand.Intn(100)),
		}

		_, err = client.
			CoreV1().
			ConfigMaps(namespace).
			Update(
				context.Background(),
				one,
				metav1.UpdateOptions{},
			)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Successfully updated ConfigMap")

		_, err = client.
			CoreV1().
			ConfigMaps(namespace).
			Update(
				context.Background(),
				two,
				metav1.UpdateOptions{},
			)
		fmt.Println(err)
		return err
	})
}

func getConfigMap(c *kubernetes.Clientset, name, ns string) (*corev1.ConfigMap, error) {
	return c.
		CoreV1().
		ConfigMaps(ns).
		Get(
			context.Background(),
			name,
			metav1.GetOptions{},
		)
}

func deleteConfigMap(c *kubernetes.Clientset, name, ns string) {
	c.CoreV1().ConfigMaps(ns).Delete(context.Background(), name, metav1.DeleteOptions{})
}
