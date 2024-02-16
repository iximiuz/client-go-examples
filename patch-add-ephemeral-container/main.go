package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	ctx       = context.Background()
	namespace = "default"
)

func main() {
	// 0. Initialize the Kubernetes client.
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

	// 1. Create a pod.
	pod, err := client.CoreV1().
		Pods(namespace).
		Create(ctx, &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "patch-add-ephemeral-container",
				Namespace: namespace,
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:    "app",
						Image:   "alpine:3",
						Command: []string{"/bin/sh", "-c", "sleep 999"},
					},
				},
			},
		}, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}
	defer func() {
		_ = client.CoreV1().Pods(pod.Namespace).Delete(ctx, pod.Name, metav1.DeleteOptions{})
	}()

	// 2. Add an ephemeral container to the pod spec.
	podWithEphemeralContainer := withDebugContainer(pod)

	// 3. Prepare the patch.
	podJSON, err := json.Marshal(pod)
	if err != nil {
		panic(err.Error())
	}

	podWithEphemeralContainerJSON, err := json.Marshal(podWithEphemeralContainer)
	if err != nil {
		panic(err.Error())
	}

	patch, err := strategicpatch.CreateTwoWayMergePatch(podJSON, podWithEphemeralContainerJSON, pod)
	if err != nil {
		panic(err.Error())
	}

	// 4. Apply the patch.
	pod, err = client.CoreV1().
		Pods(pod.Namespace).
		Patch(
			ctx,
			pod.Name,
			types.StrategicMergePatchType,
			patch,
			metav1.PatchOptions{},
			"ephemeralcontainers",
		)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Pod has %d ephemeral containers.\n", len(pod.Spec.EphemeralContainers))
}

func withDebugContainer(pod *corev1.Pod) *corev1.Pod {
	ec := &corev1.EphemeralContainer{
		EphemeralContainerCommon: corev1.EphemeralContainerCommon{
			Name:                     "debugger-123",
			Image:                    "busybox:musl",
			ImagePullPolicy:          corev1.PullIfNotPresent,
			Command:                  []string{"sh"},
			Stdin:                    true,
			TTY:                      true,
			TerminationMessagePolicy: corev1.TerminationMessageReadFile,
		},
		TargetContainerName: "app",
	}

	copied := pod.DeepCopy()
	copied.Spec.EphemeralContainers = append(copied.Spec.EphemeralContainers, *ec)

	return copied
}
