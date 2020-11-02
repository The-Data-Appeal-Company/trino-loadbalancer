package discovery

import (
	"context"
	"github.com/The-Data-Appeal-Company/presto-loadbalancer/pkg/models"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"net/url"
	"strings"
	"testing"
)

type k8sClient struct {
	kubernetes.Interface
}

type mockCoreV1 struct {
	corev1.CoreV1Interface
}

type mockNamespace struct {
	corev1.NamespaceInterface
}

type mockServiceDefault struct {
	corev1.ServiceInterface
}

type mockServiceNs1 struct {
	corev1.ServiceInterface
}

type mockServiceNs2 struct {
	corev1.ServiceInterface
}

func (c k8sClient) CoreV1() corev1.CoreV1Interface {
	return mockCoreV1{}
}

func (mc mockCoreV1) Namespaces() corev1.NamespaceInterface {
	return mockNamespace{}
}

func (mc mockCoreV1) Services(namespace string) corev1.ServiceInterface {
	switch namespace {
	case "ns-1":
		return mockServiceNs1{}
	case "ns-2":
		return mockServiceNs2{}
	}

	return mockServiceDefault{}
}

func (ms mockServiceDefault) List(ctx context.Context, opts metav1.ListOptions) (*v1.ServiceList, error) {

	return &v1.ServiceList{
		TypeMeta: metav1.TypeMeta{},
		ListMeta: metav1.ListMeta{},
		Items:    []v1.Service{},
	}, nil
}

func (ms mockServiceNs1) List(ctx context.Context, opts metav1.ListOptions) (*v1.ServiceList, error) {

	if strings.Contains(opts.LabelSelector, "presto.distribution=prestosql") {

		return &v1.ServiceList{
			TypeMeta: metav1.TypeMeta{},
			ListMeta: metav1.ListMeta{},
			Items: []v1.Service{
				{ObjectMeta: metav1.ObjectMeta{
					Name:      "prestosql-1",
					Namespace: "ns-1",
				}},
				{ObjectMeta: metav1.ObjectMeta{
					Name:      "prestosql-12",
					Namespace: "ns-1",
				}},
			},
		}, nil
	}

	return &v1.ServiceList{
		TypeMeta: metav1.TypeMeta{},
		ListMeta: metav1.ListMeta{},
		Items: []v1.Service{
			{ObjectMeta: metav1.ObjectMeta{
				Name:      "prestodb-1",
				Namespace: "ns-1",
			}},
		},
	}, nil

}

func (ms mockServiceNs2) List(ctx context.Context, opts metav1.ListOptions) (*v1.ServiceList, error) {

	if strings.Contains(opts.LabelSelector, "presto.distribution=prestosql") {
		return &v1.ServiceList{
			TypeMeta: metav1.TypeMeta{},
			ListMeta: metav1.ListMeta{},
			Items: []v1.Service{
				{ObjectMeta: metav1.ObjectMeta{
					Name:      "prestosql-2",
					Namespace: "ns-2",
				}},
			},
		}, nil
	}
	return &v1.ServiceList{
		TypeMeta: metav1.TypeMeta{},
		ListMeta: metav1.ListMeta{},
		Items:    []v1.Service{},
	}, nil
}

func (mn mockNamespace) List(ctx context.Context, opts metav1.ListOptions) (*v1.NamespaceList, error) {

	return &v1.NamespaceList{
		Items: []v1.Namespace{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "default",
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "ns-1",
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "ns-2",
				},
			},
		},
	}, nil
}

func TestK8sClusterProvider_Discover(t *testing.T) {

	clientset := fake.NewSimpleClientset()
	client := k8sClient{clientset}

	prestoUrl1, _ := url.Parse("http://prestosql-1.ns-1.svc.cluster.test")
	prestoDbUrl1, _ := url.Parse("http://prestodb-1.ns-1.svc.cluster.test")
	prestoUrl2, _ := url.Parse("http://prestosql-2.ns-2.svc.cluster.test")
	prestoUrl12, _ := url.Parse("http://prestosql-12.ns-1.svc.cluster.test")

	type fields struct {
		k8sClient     kubernetes.Interface
		SelectorTags  map[string]string
		clusterDomain string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []models.Coordinator
		wantErr bool
	}{
		{
			name: "shouldDiscoverPrestoSqlInK8sCluster",
			fields: fields{
				k8sClient: client,
				SelectorTags: map[string]string{
					"presto.distribution": "prestosql",
				},
				clusterDomain: "cluster.test",
			},
			want: []models.Coordinator{
				{
					Name: "prestosql-1",
					URL:  prestoUrl1,
					Tags: map[string]string{
						"presto.distribution": "prestosql",
					},
					Enabled:      true,
					Distribution: "prestosql",
				},
				{
					Name: "prestosql-12",
					URL:  prestoUrl12,
					Tags: map[string]string{
						"presto.distribution": "prestosql",
					},
					Enabled:      true,
					Distribution: "prestosql",
				},
				{
					Name: "prestosql-2",
					URL:  prestoUrl2,
					Tags: map[string]string{
						"presto.distribution": "prestosql",
					},
					Enabled:      true,
					Distribution: "prestosql",
				},
			},
			wantErr: false,
		},
		{
			name: "shouldDiscoverPrestoDbInK8sCluster",
			fields: fields{
				k8sClient: client,
				SelectorTags: map[string]string{
					"presto.distribution": "prestodb",
				},
				clusterDomain: "cluster.test",
			},
			want: []models.Coordinator{
				{
					Name: "prestodb-1",
					URL:  prestoDbUrl1,
					Tags: map[string]string{
						"presto.distribution": "prestodb",
					},
					Enabled:      true,
					Distribution: "prestodb",
				},
			},
			wantErr: false,
		},
		{
			name: "shouldErrorWhenUnknownDistributionType",
			fields: fields{
				k8sClient: client,
				SelectorTags: map[string]string{
					"presto.distribution": "lentodb",
				},
				clusterDomain: "cluster.test",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &K8sClusterProvider{
				k8sClient:     tt.fields.k8sClient,
				SelectorTags:  tt.fields.SelectorTags,
				clusterDomain: tt.fields.clusterDomain,
			}
			got, err := k.Discover(context.TODO())
			if (err != nil) != tt.wantErr {
				t.Errorf("Discover() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, got, tt.want)
		})
	}
}
