package main

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	namespace         = "default"
	label             = "informer-dynamic-simple-" + rand.String(6)
	ConfigMapResource = schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "configmaps",
	}
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

	// Create one object before initializing the informer.
	first := createConfigMap(client)

	// Create a shared informer factory.
	//   - A factory is essentially a struct keeping a map (type -> informer).
	//   - 5*time.Second is a default resync period (for all informers).
	factory := dynamicinformer.NewDynamicSharedInformerFactory(client, 5*time.Second)

	// When informer is requested, the factory instantiates it and keeps the
	// the reference to it in the internal map before returning.
	dynamicInformer := factory.ForResource(ConfigMapResource)
	dynamicInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			cm := obj.(*unstructured.Unstructured)
			fmt.Printf("Informer event: ConfigMap ADDED %s/%s\n", cm.GetNamespace(), cm.GetName())
		},
		UpdateFunc: func(old, new interface{}) {
			cm := old.(*unstructured.Unstructured)
			fmt.Printf("Informer event: ConfigMap UPDATED %s/%s\n", cm.GetNamespace(), cm.GetName())
		},
		DeleteFunc: func(obj interface{}) {
			cm := obj.(*unstructured.Unstructured)
			fmt.Printf("Informer event: ConfigMap DELETED %s/%s\n", cm.GetNamespace(), cm.GetName())
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the informers' machinery.
	//   - Start() starts every Informer requested before using a goroutine per informer.
	//   - A started Informer will fetch ALL the ConfigMaps from all the namespaces
	//     (using a lister) and trigger `AddFunc`` for each found ConfigMap object.
	//     Use NewSharedInformerFactoryWithOptions() to make the lister fetch only
	//     a filtered subset of objects.
	//   - All ConfigMaps added, updated, or deleted after the informer has been synced
	//     will trigger the corresponding callback call (using a watch).
	//   - Every 5*time.Second the UpdateFunc callback will be called for every
	//     previously fetched ConfigMap (so-called resync period).
	factory.Start(ctx.Done())

	// factory.Start() releases the execution flow without waiting for all the
	// internal machinery to warm up. We use factory.WaitForCacheSync() here
	// to poll for cmInformer.Informer().HasSynced(). Essentially, it's just a
	// fancy way to write a while-loop checking HasSynced() flags for all the
	// registered informers with 100ms delay between iterations.
	for gvr, ok := range factory.WaitForCacheSync(ctx.Done()) {
		if !ok {
			panic(fmt.Sprintf("Failed to sync cache for resource %v", gvr))
		}
	}

	// Search for the existing ConfigMap object using the label selector.
	selector, err := labels.Parse("example==" + label)
	if err != nil {
		panic(err.Error())
	}
	list, err := dynamicInformer.Lister().List(selector)
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

func createConfigMap(client dynamic.Interface) *unstructured.Unstructured {
	cm := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]interface{}{
				"namespace":    namespace,
				"generateName": "informer-dynamic-simple-",
				"labels": map[string]interface{}{
					"example": label,
				},
			},
			"data": map[string]interface{}{
				"foo": "bar",
			},
		},
	}

	cm, err := client.
		Resource(ConfigMapResource).
		Namespace(namespace).
		Create(context.Background(), cm, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Created ConfigMap %s/%s\n", cm.GetNamespace(), cm.GetName())
	return cm
}

func deleteConfigMap(client dynamic.Interface, cm *unstructured.Unstructured) {
	err := client.
		Resource(ConfigMapResource).
		Namespace(cm.GetNamespace()).
		Delete(context.Background(), cm.GetName(), metav1.DeleteOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Deleted ConfigMap %s/%s\n", cm.GetNamespace(), cm.GetName())
}
