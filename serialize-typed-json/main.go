package main

import (
	"encoding/json"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	jsonserializer "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
)

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

	// Typed -> JSON (Option I)
	//   - Serializer = Decoder + Encoder. Since we need only Encoder functionality
	//     in this example, we can pass nil's instead of MetaFactory, Creater, and
	//     Typer arguments as they are used only by the Decoder.
	encoder := jsonserializer.NewSerializerWithOptions(
		nil, // jsonserializer.MetaFactory
		nil, // runtime.ObjectCreater
		nil, // runtime.ObjectTyper
		jsonserializer.SerializerOptions{
			Yaml:   false,
			Pretty: false,
			Strict: false,
		},
	)

	// Runtime.Encode() is just a helper function to invoke Encoder.Encode()
	encoded, err := runtime.Encode(encoder, &obj)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Serialized (option I)", string(encoded))

	// Typed -> JSON (Option II)
	//   Actually, the implementation of Encoder.Encode() in the case of JSON
	//   boils down to calling the stdlib encoding/json.Marshal() with optional
	//   pretty-printing and converting JSON to YAML.
	//   See https://github.com/kubernetes/apimachinery/blob/73cb564852596cc976f3ead9e0f4678875af0cbf/pkg/runtime/serializer/json/json.go#L210-L234
	encoded2, err := json.Marshal(obj)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Serialized (option II)", string(encoded2))

	// JSON -> Typed
	//   - Serializer = Decoder + Encoder.
	//   - jsonserializer.MetaFactory is a simple partial JSON unmarshaller that
	//     looks for APIGroup/Version and Kind attributes in the supplied
	//     piece of JSON and parses them into a schema.GroupVersionKind{} object.
	//   - runtime.ObjectCreater is used to create an empty typed runtime.Object
	//     (e.g., Deployment, Pod, ConfigMap, etc.) for the provided APIGroup/Version and Kind.
	//   - runtime.ObjectTyper is rather optional - Decoder accepts an optional
	//     `into runtime.Object` argument, and ObjectTyper is used to make sure
	//     the MetaFactory's GroupVersionKind matches the one from the `into` argument.
	decoder := jsonserializer.NewSerializerWithOptions(
		jsonserializer.DefaultMetaFactory, // jsonserializer.MetaFactory
		scheme.Scheme,                     // runtime.Scheme implements runtime.ObjectCreater
		scheme.Scheme,                     // runtime.Scheme implements runtime.ObjectTyper
		jsonserializer.SerializerOptions{
			Yaml:   false,
			Pretty: false,
			Strict: false,
		},
	)

	// The actual decoding is much like stdlib encoding/json.Unmarshal but with some
	// minor tweaks - see https://github.com/kubernetes-sigs/json for more.
	decoded, err := runtime.Decode(decoder, encoded)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Deserialized %#v\n", decoded)
}
