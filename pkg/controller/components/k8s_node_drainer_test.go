package components

import (
	"context"
	"errors"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/logging"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestNodeDrainerDrainNodeSuccess(t *testing.T) {
	client := fake.NewSimpleClientset(&v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "trino-coordinator",
			Namespace: "trino",
		},
		Status: v1.PodStatus{
			Phase: v1.PodRunning,
		},
		Spec: v1.PodSpec{
			NodeName: "node-00",
		},
	}, &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "trino-worker-00",
			Namespace: "trino",
		},
		Status: v1.PodStatus{
			Phase: v1.PodRunning,
		},
		Spec: v1.PodSpec{
			NodeName: "node-01",
		},
	}, &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "trino",
			Labels: map[string]string{
				"scope": "trino",
			},
		},
	}, &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "node-00",
		},
	}, &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "node-01",
		},
	})

	conf := KubeNodeDrainerConf{
		namespaceLabelSelector: map[string]string{
			"scope": "trino",
		},
		podGracePeriod: 0,
	}

	nodeDrainer := NewKubeNodeDrainer(client, conf, logging.Noop())

	err := nodeDrainer.Drain(context.TODO(), "trino-worker-00")
	require.NoError(t, err)
}

func TestNodeDrainerDrainPodNotFound(t *testing.T) {
	client := fake.NewSimpleClientset(&v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "trino-coordinator",
			Namespace: "trino",
		},
		Status: v1.PodStatus{
			Phase: v1.PodRunning,
		},
		Spec: v1.PodSpec{
			NodeName: "node-00",
		},
	}, &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "trino-worker-00",
			Namespace: "trino",
		},
		Status: v1.PodStatus{
			Phase: v1.PodRunning,
		},
		Spec: v1.PodSpec{
			NodeName: "node-01",
		},
	}, &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "trino",
			Labels: map[string]string{
				"scope": "trino",
			},
		},
	})

	conf := KubeNodeDrainerConf{
		namespaceLabelSelector: map[string]string{
			"scope": "trino",
		},
		podGracePeriod: 0,
	}

	nodeDrainer := NewKubeNodeDrainer(client, conf, logging.Noop())

	err := nodeDrainer.Drain(context.TODO(), "trino-worker-01")
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrPodNotFound))
}

func TestNodeDrainerDrainNodeNotFound(t *testing.T) {
	client := fake.NewSimpleClientset(&v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "trino-coordinator",
			Namespace: "trino",
		},
		Status: v1.PodStatus{
			Phase: v1.PodRunning,
		},
		Spec: v1.PodSpec{
			NodeName: "node-00",
		},
	}, &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "trino-worker-00",
			Namespace: "trino",
		},
		Status: v1.PodStatus{
			Phase: v1.PodRunning,
		},
		Spec: v1.PodSpec{
			NodeName: "node-01",
		},
	}, &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "trino",
			Labels: map[string]string{
				"scope": "trino",
			},
		},
	}, &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "node-00",
		},
	})

	conf := KubeNodeDrainerConf{
		namespaceLabelSelector: map[string]string{
			"scope": "trino",
		},
		podGracePeriod: 0,
	}

	nodeDrainer := NewKubeNodeDrainer(client, conf, logging.Noop())

	err := nodeDrainer.Drain(context.TODO(), "trino-worker-00")
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrNodeNotFound))
}

func TestNodeDrainerDrainPodNotInNamespaceSelector(t *testing.T) {
	client := fake.NewSimpleClientset(&v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "trino-coordinator",
			Namespace: "trino",
		},
		Status: v1.PodStatus{
			Phase: v1.PodRunning,
		},
		Spec: v1.PodSpec{
			NodeName: "node-00",
		},
	}, &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "trino-worker-00",
			Namespace: "trino",
		},
		Status: v1.PodStatus{
			Phase: v1.PodRunning,
		},
		Spec: v1.PodSpec{
			NodeName: "node-01",
		},
	}, &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "trino",
			Labels: map[string]string{
				"scope": "trino",
			},
		},
	}, &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "node-00",
		},
	}, &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "node-01",
		},
	})

	conf := KubeNodeDrainerConf{
		namespaceLabelSelector: map[string]string{
			"scope": "trino-other",
		},
		podGracePeriod: 0,
	}

	nodeDrainer := NewKubeNodeDrainer(client, conf, logging.Noop())

	err := nodeDrainer.Drain(context.TODO(), "trino-worker-00")
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrPodNotFound))
}

func TestNodeDrainerMultiNamespace(t *testing.T) {
	client := fake.NewSimpleClientset(&v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "trino-worker-00",
			Namespace: "trino-00",
		},
		Status: v1.PodStatus{
			Phase: v1.PodRunning,
		},
		Spec: v1.PodSpec{
			NodeName: "node-00",
		},
	}, &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "trino-worker-01",
			Namespace: "trino-01",
		},
		Status: v1.PodStatus{
			Phase: v1.PodRunning,
		},
		Spec: v1.PodSpec{
			NodeName: "node-01",
		},
	}, &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "trino-01",
			Labels: map[string]string{
				"scope": "trino",
			},
		},
	},&v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "trino-00",
			Labels: map[string]string{
				"scope": "trino",
			},
		},
	}, &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "node-00",
		},
	}, &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "node-01",
		},
	})

	conf := KubeNodeDrainerConf{
		namespaceLabelSelector: map[string]string{
			"scope": "trino",
		},
		podGracePeriod: 0,
	}

	nodeDrainer := NewKubeNodeDrainer(client, conf, logging.Noop())

	err := nodeDrainer.Drain(context.TODO(), "trino-worker-01")
	require.NoError(t, err)
}
