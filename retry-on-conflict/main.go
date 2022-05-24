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
	defer deleteConfigMap(client, name, namespace)

	firstTry := true
	err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// always fetch the new version from the api server
		c, err := getConfigMap(client, name, namespace)
		if err != nil {
			return err
		}

		// simulate an external update
		// this code will not exist in a typical implementation
		// this is just for demonstration
		if firstTry {
			firstTry = false
			simulateExternalUpdate(client, name, namespace)
		}

		// generate random int so we're always different
		c.Data = map[string]string{
			"changed": fmt.Sprint(rand.Intn(100)),
		}

		_, err = client.
			CoreV1().
			ConfigMaps(namespace).
			Update(
				context.Background(),
				c,
				metav1.UpdateOptions{},
			)
		if err != nil {
			fmt.Println(err)
		}
		return err
	})
	if err != nil {
		// ensure no other error type occurred
		panic(err)
	}
	fmt.Println("Successfully updated ConfigMap")
}

func simulateExternalUpdate(k *kubernetes.Clientset, name, ns string) {
	cm, err := getConfigMap(k, name, ns)
	if err != nil {
		panic(err)
	}
	cm.Data = map[string]string{
		"external": "update",
	}
	_, err = k.CoreV1().
		ConfigMaps(ns).
		Update(
			context.Background(),
			cm,
			metav1.UpdateOptions{},
		)
	if err != nil {
		panic(err)
	}
}

func getConfigMap(k *kubernetes.Clientset, name, ns string) (*corev1.ConfigMap, error) {
	return k.
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
