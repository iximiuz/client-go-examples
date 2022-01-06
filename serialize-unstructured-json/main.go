package main

import (
	"fmt"
	"reflect"

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

	// Unstructured -> JSON (Option I)
	//   - Despite the name, `UnstructuredJSONScheme` is not a scheme but a codec
	//   - runtime.Encode() is just a helper function to invoke UnstructuredJSONScheme.Encode()
	//   - UnstructuredJSONScheme.Encode() is needed because the unstructured instance can be
	//     either a single object, a list, or an unknown runtime object, so some amount of
	//     preprocessing is required before passing the data to json.Marshal()
	//   - Usage example: dynamic client (client-go/dynamic.Interface)
	bytes, err := runtime.Encode(unstructured.UnstructuredJSONScheme, &uConfigMap)
	fmt.Println("Serialized (option I)", string(bytes))

	// Unstructured -> JSON (Option II)
	//   - This is just a handy shortcut for the above code.
	bytes, err = uConfigMap.MarshalJSON()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Serialized (option II)", string(bytes))

	// JSON -> Unstructured (Option I)
	//   - Usage example: dynamic client (client-go/dynamic.Interface)
	obj1, err := runtime.Decode(unstructured.UnstructuredJSONScheme, bytes)
	if err != nil {
		panic(err.Error())
	}

	// JSON -> Unstructured (Option II)
	//   - This is just a handy shortcut for the above code.
	obj2 := &unstructured.Unstructured{}
	err = obj2.UnmarshalJSON(bytes)
	if err != nil {
		panic(err.Error())
	}
	if !reflect.DeepEqual(obj1, obj2) {
		panic("Unexpected configmap data")
	}
}
