package main

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
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

	// Create one object before initializing the informer.
	first := createConfigMap(client)

	// Create a shared ConfigMap informer using the factory.
	// 5*time.Second is a default resync period (for all informers).
	factory := informers.NewSharedInformerFactory(client, 5*time.Second)
	cmInformer := factory.Core().V1().ConfigMaps()
	cmInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			cm := obj.(*corev1.ConfigMap)
			fmt.Printf("Informer event: ConfigMap ADDED %s/%s\n", cm.GetNamespace(), cm.GetName())
		},
		UpdateFunc: func(old, new interface{}) {
			cm := old.(*corev1.ConfigMap)
			fmt.Printf("Informer event: ConfigMap UPDATED %s/%s\n", cm.GetNamespace(), cm.GetName())
		},
		DeleteFunc: func(obj interface{}) {
			cm := obj.(*corev1.ConfigMap)
			fmt.Printf("Informer event: ConfigMap DELETED %s/%s\n", cm.GetNamespace(), cm.GetName())
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the informer machinery.
	//   - Informer will fetch ALL the ConfigMaps from all the namespaces and trigger
	//     AddFunc for each found ConfigMap object.
	//     Use NewFilteredSharedInformerFactory to fetch only a filtered subset of objects.
	//   - All ConfigMaps added, updated, or deleted after the informer is synced
	//     will trigger the corresponding callback call.
	//   - Every 5*time.Second the UpdateFunc callback will be called for every
	//     previously fetched ConfigMap (so-called resync period).
	factory.Start(ctx.Done())
	if !cache.WaitForNamedCacheSync("my-example", ctx.Done(), cmInformer.Informer().HasSynced) {
		panic("Failed to sync cache")
	}

	// Search for the existing ConfigMap object using the label selector.
	selector, err := labels.Parse("example==" + label)
	if err != nil {
		panic(err.Error())
	}
	list, err := cmInformer.Lister().List(selector)
	if err != nil {
		panic(err.Error())
	}
	if len(list) != 1 {
		panic("expected ConfigMap not found")
	}

	// Create another object while watching.
	second := createConfigMap(client)

	// Delete config maps created by this test.
	deleteConfigMap(client, first)
	deleteConfigMap(client, second)

	// Stay for a couple more seconds to observe resyncs.
	time.Sleep(10 * time.Second)
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
