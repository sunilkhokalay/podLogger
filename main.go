package main

import (
	"context"
	"flag"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

func main() {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		var kubeconfig *string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()

		// use the current context in kubeconfig
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	}

	// creates the clientSet
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	wInf, err := clientSet.CoreV1().Pods("").Watch(context.TODO(), metav1.ListOptions{
		TypeMeta: metav1.TypeMeta{},
		Watch:    true,
	})

	if err != nil {
		panic(err)
	}

	eventChan := wInf.ResultChan()
	println("Starting pod watcher...")
	for {
		event := <-eventChan
		pod := event.Object.(*v1.Pod)
		println(fmt.Sprintf("The pod \"%s\" in \"%s\" namespace was %s", pod.Name, pod.Namespace, event.Type))
	}
}
