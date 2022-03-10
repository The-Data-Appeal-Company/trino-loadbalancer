package autoscaler

import (
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/trino"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/logging"
	testUtil "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/tests"
	"k8s.io/client-go/kubernetes"
	"reflect"
	"testing"
	"time"
)

func Test_lastQueryExecution(t *testing.T) {
	type args struct {
		queries trino.QueryList
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "test get last execution from empty array",
			args: args{},
			want: time.Time{},
		},
		{
			name: "test get last execution from queries",
			args: args{
				queries: []trino.QueryListItem{
					{
						QueryStats: trino.QueryStats{
							EndTime: time.Now(),
						},
					},
					{
						QueryStats: trino.QueryStats{
							EndTime: time.Unix(1635956820, 0),
						},
					},
					{
						QueryStats: trino.QueryStats{
							EndTime: time.Unix(1646843221, 0),
						},
					},
				},
			},
			want: time.Unix(1635956820, 0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lastQueryExecution(tt.args.queries); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("lastQueryExecution() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hasQueriesInState(t *testing.T) {
	type args struct {
		queries trino.QueryList
		state   string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should return false on empty slice",
			args: args{},
			want: false,
		},
		{
			name: "test return true with atleast one query in the selected state",
			args: args{
				queries: []trino.QueryListItem{
					{
						State: StateRunning,
					},
					{
						State: StateWaitingForResources,
					},
				},
				state: StateRunning,
			},
			want: true,
		},
		{
			name: "test return true with atleast one query in the selected state",
			args: args{
				queries: []trino.QueryListItem{
					{
						State: StateRunning,
					},
					{
						State: StateWaitingForResources,
					},
				},
				state: StateWaitingForResources,
			},
			want: true,
		},
		{
			name: "test return false when no query is in the selected state",
			args: args{
				queries: []trino.QueryListItem{
					{
						State: StateWaitingForResources,
					},
					{
						State: StateWaitingForResources,
					},
				},
				state: StateRunning,
			},
			want: false,
		},
		{
			name: "test return true when all queries are in the selected state",
			args: args{
				queries: []trino.QueryListItem{
					{
						State: StateRunning,
					},
					{
						State: StateRunning,
					},
				},
				state: StateRunning,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasQueriesInState(tt.args.queries, tt.args.state); got != tt.want {
				t.Errorf("hasQueriesInState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKubeClientAutoscaler_needScaleDown(t *testing.T) {
	type fields struct {
		client   kubernetes.Interface
		trinoApi trino.Api
		state    State
	}
	type args struct {
		req     KubeRequest
		queries trino.QueryList
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:   "don't scale down when no queries and no state are present",
			fields: fields{state: MemoryState()},
			args: args{
				req: KubeRequest{
					Coordinator: testUtil.MustUrl("http://coordinator.local"),
					ScaleAfter:  10 * time.Second,
				},
				queries: nil,
			},
			want:    false,
			wantErr: false,
		},
		{
			name:   "scale down when the last query has been executed after now-'ScaleAfter'",
			fields: fields{state: MemoryState()},
			args: args{
				req: KubeRequest{
					Coordinator: testUtil.MustUrl("http://coordinator.local"),
					ScaleAfter:  30 * time.Minute,
				},
				queries: trino.QueryList{
					{
						State: "COMPLETED",
						QueryStats: trino.QueryStats{
							EndTime: time.Now().Add(-1 * time.Hour),
						},
					},
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name:   "dont scale down when with running queries",
			fields: fields{state: MemoryState()},
			args: args{
				req: KubeRequest{
					Coordinator: testUtil.MustUrl("http://coordinator.local"),
					ScaleAfter:  30 * time.Minute,
				},
				queries: trino.QueryList{
					{
						State: "COMPLETED",
						QueryStats: trino.QueryStats{
							EndTime: time.Now().Add(-1 * time.Hour),
						},
					},
					{
						State: StateRunning,
					},
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "scale down when no queries are returned but state is present",
			fields: fields{state: mockState{
				set: func(clusterID string, t time.Time) error {
					return nil
				},
				get: func(clusterID string) (time.Time, error) {
					return time.Now().Add(-1 * time.Hour), nil
				},
			}},
			args: args{
				req: KubeRequest{
					Coordinator: testUtil.MustUrl("http://coordinator.local"),
					ScaleAfter:  30 * time.Minute,
				},
				queries: trino.QueryList{},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &KubeClientAutoscaler{
				client:   tt.fields.client,
				trinoApi: tt.fields.trinoApi,
				state:    tt.fields.state,
				logger:   logging.Noop(),
			}
			got, err := k.needScaleDown(tt.args.req, tt.args.queries)
			if (err != nil) != tt.wantErr {
				t.Errorf("needScaleDown() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("needScaleDown() got = %v, want %v", got, tt.want)
			}
		})
	}
}
