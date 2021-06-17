package components

import (
	"context"
	"errors"
	"fmt"
	ioc "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/io"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/logging"
	v12 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubectl/pkg/drain"
	"time"
)

var (
	ErrPodNotFound  = errors.New("specified pod not found")
	ErrNodeNotFound = errors.New("node not found")
)

type NodeDrainer interface {
	Drain(ctx context.Context, nodeID string) error
}

type KubeNodeDrainerConf struct {
	namespaceLabelSelector map[string]string
	podGracePeriod         time.Duration
}

type KubeNodeDrainer struct {
	client kubernetes.Interface
	conf   KubeNodeDrainerConf
	logger logging.Logger
}

func NewKubeNodeDrainer(client kubernetes.Interface, conf KubeNodeDrainerConf, logger logging.Logger) KubeNodeDrainer {
	return KubeNodeDrainer{
		client: client,
		conf:   conf,
		logger: logger,
	}
}

func (k KubeNodeDrainer) Drain(ctx context.Context, nodeID string) error {
	// the nodeID is actually the pod name, we don't know in which namespace the pod is located, so we have to perform
	// a list operation in order to have the Pod reference
	var podName = nodeID

	namespaces, err := k.client.CoreV1().Namespaces().List(ctx, v1.ListOptions{
		LabelSelector: labels.FormatLabels(k.conf.namespaceLabelSelector),
	})

	if err != nil {
		return err
	}

	for _, ns := range namespaces.Items {

		pods, err := k.client.CoreV1().Pods(ns.Name).List(ctx, v1.ListOptions{})
		if err != nil {
			return err
		}

		for _, pod := range pods.Items {
			if pod.Status.Phase != v12.PodRunning {
				continue
			}

			if pod.Name != podName {
				continue
			}

			var nodeName = pod.Spec.NodeName
			if len(nodeName) == 0 {
				return fmt.Errorf("no node found for pod %s/%s", ns.Name, podName)
			}

			k.logger.Info("associated pod %s with node %s", podName, nodeName)

			if err := k.cordonAndDrainNode(ctx, nodeName); err != nil {
				if kerrors.IsNotFound(err) {
					return ErrNodeNotFound
				}
				return err
			}

			return nil
		}

	}

	return ErrPodNotFound
}

func (k KubeNodeDrainer) cordonAndDrainNode(ctx context.Context, nodeName string) error {
	node, err := k.client.CoreV1().Nodes().Get(ctx, nodeName, v1.GetOptions{})
	if err != nil {
		return err
	}

	var drainHelper = &drain.Helper{
		Ctx:                 ctx,
		Client:              k.client,
		Force:               true,
		Timeout:             k.conf.podGracePeriod * time.Second,
		GracePeriodSeconds:  int(k.conf.podGracePeriod.Seconds()),
		IgnoreAllDaemonSets: true,
		DeleteEmptyDirData:  true,
		ErrOut:              ioc.NewNoopWriter(),
		Out:                 ioc.NewNoopWriter(),
	}

	if err := drain.RunCordonOrUncordon(drainHelper, node, true); err != nil {
		return err
	}

	if err := drain.RunNodeDrain(drainHelper, nodeName); err != nil {
		return err
	}

	return nil
}
