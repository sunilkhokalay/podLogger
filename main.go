package main

import (
	"context"
	"flag"
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

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	wInf, err := clientset.CoreV1().Pods("default").Watch(context.TODO(),metav1.ListOptions{
		TypeMeta:            metav1.TypeMeta{},
		Watch:               true,
	})

	if err != nil{
		panic(err)
	}

	eventChan := wInf.ResultChan()

	for {
		event := <- eventChan
		println("---------------")
		println("Event received:-")
		println(event.Type)
		println(event.Object.GetObjectKind().GroupVersionKind().Group)
		println(event.Object.GetObjectKind().GroupVersionKind().Version)
		println(event.Object.GetObjectKind().GroupVersionKind().Kind)
		println("---------------")
	}
}
