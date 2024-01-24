package kubernetes

import (
	"flag"
	"herald/env"
	"log"
	"os/user"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	clientset kubernetes.Interface
}

func NewClient() (*Client, error) {
	config, err := getConfig()
	if err != nil {
		return nil, err
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		clientset: clientSet,
	}, nil
}

func getConfig() (config *rest.Config, err error) {
	if env.IsInCluster() {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
		log.Printf("using inClusterConfig")
	} else {
		usr, err := user.Current()
		if err != nil {
			return nil, err
		}

		filePath := usr.HomeDir + "/.kube/config"
		kubeconfig := flag.String("kubeconfig", filePath, "absolute path to file")
		flag.Parse()
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}
