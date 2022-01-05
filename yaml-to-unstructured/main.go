package main

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func main() {
	yConfigMap := `---
apiVersion: v1
data:
  foo: bar
kind: ConfigMap
metadata:
  creationTimestamp:
  name: my-configmap
  namespace: default
`

	// YAML -> Unstructured (through JSON)
	jConfigMap, err := yaml.ToJSON([]byte(yConfigMap))
	if err != nil {
		panic(err.Error())
	}

	object, err := runtime.Decode(unstructured.UnstructuredJSONScheme, jConfigMap)
	if err != nil {
		panic(err.Error())
	}

	uConfigMap, ok := object.(*unstructured.Unstructured)
	if !ok {
		panic("unstructured.Unstructured expected")
	}

	if uConfigMap.GetName() != "my-configmap" {
		panic("Unexpected configmap data")
	}
}
