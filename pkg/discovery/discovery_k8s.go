package discovery

import (
	"context"
	"fmt"
	"github.com/The-Data-Appeal-Company/presto-loadbalancer/pkg/models"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	k8s "k8s.io/client-go/kubernetes"
	"net/url"
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

		logrus.Infof("namespace %s", ns.Name)

		services, err := k.k8sClient.CoreV1().Services(ns.Name).List(ctx, v1.ListOptions{
			LabelSelector: labels.FormatLabels(k.SelectorTags),
		})

		if err != nil {
			return nil, err
		}

		for _, svc := range services.Items {

			logrus.Infof("service %s", svc.Name)

			svcUrl, err := url.Parse(fmt.Sprintf("http://%s.%s.svc.%s", svc.Name, svc.Namespace, k.clusterDomain))
			if err != nil {
				return nil, err
			}

			dist, err := models.ParsePrestoDist(k.SelectorTags["presto.distribution"])

			if err != nil {
				return nil, err
			}

			coordinators = append(coordinators, models.Coordinator{
				Name:         svc.Name,
				URL:          svcUrl,
				Tags:         k.SelectorTags,
				Enabled:      true,
				Distribution: dist,
			})
		}
	}

	return coordinators, nil

}
