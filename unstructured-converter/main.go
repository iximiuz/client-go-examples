package main

import (
	"reflect"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func main() {
	uConfigMap := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]interface{}{
				"creationTimestamp": nil,
				"namespace":         "default",
				"name":              "my-configmap",
			},
			"data": map[string]interface{}{
				"foo": "bar",
			},
		},
	}

	// Unstructured -> Typed
	var tConfigMap corev1.ConfigMap
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(uConfigMap.Object, &tConfigMap)
	if err != nil {
		panic(err.Error())
	}
	if tConfigMap.GetName() != "my-configmap" {
		panic("Typed config map has unexpected data")
	}

	// Typed -> Unstructured
	object, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&tConfigMap)
	if err != nil {
		panic(err.Error())
	}
	if !reflect.DeepEqual(unstructured.Unstructured{Object: object}, uConfigMap) {
		panic("Unstructured config map has unexpected data")
	}
}
