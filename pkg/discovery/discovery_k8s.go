package discovery

import (
	"context"
	"fmt"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
	"github.com/sirupsen/logrus"
	v12 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	k8s "k8s.io/client-go/kubernetes"
	"net/url"
)

const (
	svcPortName = "http"
)

type K8sClusterProvider struct {
	k8sClient     k8s.Interface
	SelectorTags  map[string]string
	clusterDomain string
}

func NewK8sClusterProvider(k8sClient k8s.Interface, selectorTags map[string]string, clusterDomain string) *K8sClusterProvider {
	return &K8sClusterProvider{k8sClient: k8sClient, SelectorTags: selectorTags, clusterDomain: clusterDomain}
}

func (k *K8sClusterProvider) Discover(ctx context.Context) ([]models.Coordinator, error) {

	coordinators := make([]models.Coordinator, 0)
	namespaces, err := k.k8sClient.CoreV1().Namespaces().List(ctx, v1.ListOptions{})

	if err != nil {
		return nil, err
	}

	for _, ns := range namespaces.Items {
		services, err := k.k8sClient.CoreV1().Services(ns.Name).List(ctx, v1.ListOptions{
			LabelSelector: labels.FormatLabels(k.SelectorTags),
		})

		if err != nil {
			return nil, err
		}

		for _, svc := range services.Items {

			logrus.Infof("service %s", svc.Name)

			servicePort, err := portByName(svc.Spec.Ports, svcPortName)
			if err != nil {
				return nil, err
			}

			svcUrl, err := url.Parse(fmt.Sprintf("http://%s.%s.svc.%s:%d", svc.Name, svc.Namespace, k.clusterDomain, servicePort.Port))
			if err != nil {
				return nil, err
			}

			coordinators = append(coordinators, models.Coordinator{
				Name:    fmt.Sprintf("%s-%s", svc.Namespace, svc.Name),
				URL:     svcUrl,
				Tags:    k.SelectorTags,
				Enabled: true,
			})
		}
	}

	return coordinators, nil

}

func portByName(ports []v12.ServicePort, name string) (v12.ServicePort, error) {
	for _, port := range ports {
		if port.Name == name {
			return port, nil
		}
	}

	return v12.ServicePort{}, fmt.Errorf("no port with name %s found", name)
}
