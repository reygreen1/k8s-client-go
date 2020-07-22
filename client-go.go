package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	ns := "k8s-cicd"
	fmt.Println("namespace: ", ns)

	clientSet := initClientSet(ns)
	// initNamespace(ns, clientSet)
	// createPod(ns, clientSet)
	// listPods(ns, clientSet)

	// time.Sleep(10 * time.Second)
	deleteNamespace(ns, clientSet)
}

func initClientSet(ns string) *kubernetes.Clientset {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	// uses the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

func initNamespace(ns string, clientset *kubernetes.Clientset) {
	namespace := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: ns,
		},
	}
	// creates namespace
	_, err := clientset.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("created namespace: ", ns)
}

func deleteNamespace(ns string, clientset *kubernetes.Clientset) {
	os.RemoveAll(
		filepath.Join(
			"/tmp",
			"drone",
			ns,
		),
	)

	// delete namespace
	clientset.CoreV1().Namespaces().Delete(
		context.TODO(),
		ns,
		metav1.DeleteOptions{},
	)
	fmt.Println("delete namespace: ", ns)
}

func createPod(ns string, clientset *kubernetes.Clientset) {
	// create pods
	pod := toPod(ns)
	_, err2 := clientset.CoreV1().Pods(ns).Create(context.TODO(), pod, metav1.CreateOptions{})
	if err2 != nil {
		panic(err2.Error())
	}
	fmt.Println("pod created!")
}

func listPods(ns string, clientset *kubernetes.Clientset) {
	pods, err := clientset.CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func toPod(namespace string) *v1.Pod {
	var volumes []v1.Volume
	volumes = append(volumes, toVolumes(namespace)...)
	var mounts []v1.VolumeMount
	mounts = append(mounts, toVolumeMounts(namespace)...)

	return &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "k8s-cicd-pod",
			Namespace: namespace,
			Labels: map[string]string{
				"app": "k8s-cicd-pod",
			},
		},
		Spec: v1.PodSpec{
			AutomountServiceAccountToken: boolptr(false),
			RestartPolicy:                v1.RestartPolicyNever,
			Containers: []v1.Container{{
				Name:         "k8s-cicd-container",
				Image:        "node:8.9.4",
				Command:      []string{"npm run build"},
				WorkingDir:   "",
				VolumeMounts: mounts,
			}},
			Volumes: volumes,
		},
	}
}

func toVolumes(namespace string) []v1.Volume {
	var to []v1.Volume
	volName := "k8s-vol"
	path := "/k8s-vol"
	volume := v1.Volume{Name: volName}
	source := v1.HostPathDirectoryOrCreate

	// NOTE the empty_dir cannot be shared across multiple
	// pods so we emulate its behavior, and mount a temp
	// directory on the host machine that can be shared
	// between pods. This means we are responsible for deleting
	// these directories.
	volume.HostPath = &v1.HostPathVolumeSource{
		Path: filepath.Join("/tmp", "drone", namespace, path),
		Type: &source,
	}
	to = append(to, volume)
	return to
}

func toVolumeMounts(namespace string) []v1.VolumeMount {
	var to []v1.VolumeMount
	volName := "k8s-vol"
	path := "/k8s-vol"
	to = append(to, v1.VolumeMount{
		Name:      volName,
		MountPath: path,
	})
	return to
}

func boolptr(v bool) *bool {
	return &v
}
