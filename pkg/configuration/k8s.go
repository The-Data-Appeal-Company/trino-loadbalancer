package configuration

import (
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

func newConfiguration(kubeConfig *string) (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err == nil {
		return config, nil
	}

	// we are not in the k8s cluster, try to use local NewConfiguration
	if kubeConfig == nil || *kubeConfig == "" {
		homeKubeConf := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		kubeConfig = &homeKubeConf
	}

	return clientcmd.BuildConfigFromFlags("", *kubeConfig)
}

func NewK8sClient(kubeConfig *string) (k8s.Interface, error) {
	kubeConf, err := newConfiguration(kubeConfig)
	if err != nil {
		return nil, err
	}

	client, err := k8s.NewForConfig(kubeConf)
	if err != nil {
		return nil, err
	}

	return client, nil
}
