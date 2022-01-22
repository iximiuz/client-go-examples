package main

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	jsonserializer "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
)

// See serialize-typed-json example for more detailed explanation
// of encoding/decoding machinery.

func main() {
	obj := corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		Data: map[string]string{"foo": "bar"},
	}
	obj.Namespace = "default"
	obj.Name = "my-configmap"

	// Serializer = Decoder + Encoder.
	serializer := jsonserializer.NewSerializerWithOptions(
		jsonserializer.DefaultMetaFactory, // jsonserializer.MetaFactory
		scheme.Scheme,                     // runtime.Scheme implements runtime.ObjectCreater
		scheme.Scheme,                     // runtime.Scheme implements runtime.ObjectTyper
		jsonserializer.SerializerOptions{
			Yaml:   true,
			Pretty: false,
			Strict: false,
		},
	)

	// Typed -> YAML
	// Runtime.Encode() is just a helper function to invoke Encoder.Encode()
	yaml, err := runtime.Encode(serializer, &obj)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Serialized:\n%s", string(yaml))

	// YAML -> Typed (through JSON, actually)
	decoded, err := runtime.Decode(serializer, yaml)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Deserialized: %#v\n", decoded)
}
