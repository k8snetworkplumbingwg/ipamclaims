package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/golang/glog"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/k8snetworkplumbingwg/ipamclaims/pkg/crd/ipamclaims/v1alpha1"
	clientset "github.com/k8snetworkplumbingwg/ipamclaims/pkg/crd/ipamclaims/v1alpha1/apis/clientset/versioned"
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

	// create a persistent IP allocation
	pip := &v1alpha1.IPAMClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: "example",
		},
		Spec: v1alpha1.IPAMClaimSpec{
			Network:   "tenantblue",
			Interface: "iface321",
		},
	}

	ipamClaim, err := exampleClient.K8sV1alpha1().IPAMClaims("default").Create(
		context.Background(),
		pip,
		metav1.CreateOptions{},
	)
	if err != nil {
		glog.Fatalf("Error creating a dummy persistentIP object: %v", err)
	}

	defer func() {
		// teardown persistent IP
		_ = exampleClient.K8sV1alpha1().IPAMClaims("default").Delete(
			context.Background(),
			pip.Name,
			metav1.DeleteOptions{},
		)
	}()

	ipamClaim.Status.IPs = []v1alpha1.CIDR{"winner", "winner", "chicken", "dinner"}
	_, err = exampleClient.K8sV1alpha1().IPAMClaims("default").UpdateStatus(
		context.Background(),
		ipamClaim,
		metav1.UpdateOptions{},
	)
	if err != nil {
		glog.Fatalf("Error creating a dummy persistentIP object: %v", err)
	}

	allPersistentIPs, err := exampleClient.K8sV1alpha1().IPAMClaims(metav1.NamespaceAll).List(
		context.Background(),
		metav1.ListOptions{},
	)
	if err != nil {
		glog.Fatalf("Error listing all persistentIP objects: %v", err)
	}

	for _, persistentIP := range allPersistentIPs.Items {
		fmt.Printf("IPAM claim name: %q\n", persistentIP.Name)
		fmt.Printf("  - spec: %v\n", persistentIP.Spec)
		fmt.Printf("  - status: %v\n", persistentIP.Status)
	}
}
