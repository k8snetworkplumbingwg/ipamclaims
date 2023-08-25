package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/golang/glog"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"

	clientset "github.com/maiqueb/persistentips/pkg/crd/persistentip/v1alpha1/apis/clientset/versioned"
)

var (
	kubeConfig = flag.String("kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	kubeAPIURL = flag.String("kube-api-url", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
)

func main() {
	flag.Parse()

	cfg, err := clientcmd.BuildConfigFromFlags(*kubeAPIURL, *kubeConfig)
	if err != nil {
		glog.Fatalf("Error building kubeconfig: %v", err)
	}

	exampleClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building example clientset: %v", err)
	}

	allPersistentIPs, err := exampleClient.K8sV1alpha1().IPAMLeases(metav1.NamespaceAll).List(
		context.Background(),
		metav1.ListOptions{},
	)
	if err != nil {
		glog.Fatalf("Error listing all persistentIP objects: %v", err)
	}

	for _, persistentIP := range allPersistentIPs.Items {
		fmt.Printf("IPAM lease name: %q\n", persistentIP.Name)
		fmt.Printf("  - spec: %v\n", persistentIP.Spec)
	}
}
