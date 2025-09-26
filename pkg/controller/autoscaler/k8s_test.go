package autoscaler

import (
	"fmt"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/trino"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/logging"
	testUtil "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/tests"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/configuration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
							EndTime: time.Unix(1735956820, 0),
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
			want: time.Unix(1735956820, 0),
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
		states  []string
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
				states: []string{StateRunning},
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
				states: []string{StateWaitingForResources},
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
				states: []string{StateRunning},
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
				states: []string{StateRunning},
			},
			want: true,
		},
		{
			name: "test return true when all queries are in one of the multiple state",
			args: args{
				queries: []trino.QueryListItem{
					{
						State: StateRunning,
					},
					{
						State: StateRunning,
					},
				},
				states: []string{StateWaitingForResources, StateRunning},
			},
			want: true,
		},
		{
			name: "test return false when all queries are in no-one of the multiple state",
			args: args{
				queries: []trino.QueryListItem{
					{
						State: StateRunning,
					},
					{
						State: StateRunning,
					},
				},
				states: []string{"COMPLETED", "FAILED"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasQueriesInStates(tt.args.queries, tt.args.states); got != tt.want {
				t.Errorf("hasQueriesInState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKubeClientAutoscaler_needScaleToZero(t *testing.T) {
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
				setTime: func(clusterID string, t time.Time) error {
					return nil
				},
				getTime: func(clusterID string) (time.Time, error) {
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
			got, err := k.needScaleToZero(tt.args.req, tt.args.queries)
			if (err != nil) != tt.wantErr {
				t.Errorf("needScaleToZero() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("needScaleToZero() got = %v, want %v", got, tt.want)
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
		current int
		wanted  int
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantScale bool
		wantInst  int
		wantErr   bool
	}{
		{
			name:   "don't scale down when no queries are running",
			fields: fields{state: MemoryState()},
			args: args{
				req: KubeRequest{
					Coordinator: testUtil.MustUrl("http://coordinator.local"),
					ScaleAfter:  10 * time.Second,
				},
				queries: nil,
			},
			wantScale: false,
			wantInst:  0,
			wantErr:   false,
		},
		{
			name: "don't scale down when queries are running but not state",
			fields: fields{state: mockState{
				setLastScaleUp: func(clusterID string, i int32, tim time.Time) error {
					assert.Equal(t, int32(20), i)
					return nil
				},
				getLastScaleUp: func(clusterID string) (int32, time.Time, error) {
					return 0, time.Now(), NoLastScaleUpStateError
				},
			}},
			args: args{
				req: KubeRequest{
					Coordinator: testUtil.MustUrl("http://coordinator.local"),
					ScaleAfter:  10 * time.Second,
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
				current: 20,
				wanted:  10,
			},
			wantScale: false,
			wantInst:  0,
			wantErr:   false,
		},
		{
			name: "Scale same instances update state",
			fields: fields{state: mockState{
				setLastScaleUp: func(clusterID string, i int32, tim time.Time) error {
					assert.Equal(t, int32(10), i)
					return nil
				},
				getLastScaleUp: func(clusterID string) (int32, time.Time, error) {
					assert.Fail(t, "should not be called")
					return 0, time.Now(), NoLastScaleUpStateError
				},
			}},
			args: args{
				req: KubeRequest{
					Coordinator: testUtil.MustUrl("http://coordinator.local"),
					ScaleAfter:  10 * time.Second,
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
				current: 10,
				wanted:  10,
			},
			wantScale: false,
			wantInst:  0,
			wantErr:   false,
		},
		{
			name: "don't scale down instaces on state is different from current",
			fields: fields{state: mockState{
				setLastScaleUp: func(clusterID string, i int32, tim time.Time) error {
					assert.Equal(t, int32(20), i)
					return nil
				},
				getLastScaleUp: func(clusterID string) (int32, time.Time, error) {
					return 10, time.Now(), nil
				},
			}},
			args: args{
				req: KubeRequest{
					Coordinator: testUtil.MustUrl("http://coordinator.local"),
					ScaleAfter:  10 * time.Second,
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
				current: 20,
				wanted:  10,
			},
			wantScale: false,
			wantInst:  0,
			wantErr:   false,
		},
		{
			name: "Scale down to less instances elapsed more than ScaleAfter from last scale up",
			fields: fields{state: mockState{
				setLastScaleUp: func(clusterID string, i int32, tim time.Time) error {
					assert.Fail(t, "should not be called")
					return nil
				},
				getLastScaleUp: func(clusterID string) (int32, time.Time, error) {
					return 20, time.Now().Add(-2 * time.Hour), nil
				},
			}},
			args: args{
				req: KubeRequest{
					Coordinator: testUtil.MustUrl("http://coordinator.local"),
					ScaleAfter:  10 * time.Second,
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
				current: 20,
				wanted:  10,
			},
			wantScale: true,
			wantInst:  10,
			wantErr:   false,
		},
		{
			name: "don't scale down to less instances not elapsed more than ScaleAfter from last scale up",
			fields: fields{state: mockState{
				setLastScaleUp: func(clusterID string, i int32, tim time.Time) error {
					assert.Fail(t, "should not be called")
					return nil
				},
				getLastScaleUp: func(clusterID string) (int32, time.Time, error) {
					return 20, time.Now().Add(-1 * time.Hour), nil
				},
			}},
			args: args{
				req: KubeRequest{
					Coordinator: testUtil.MustUrl("http://coordinator.local"),
					ScaleAfter:  5 * time.Hour,
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
				current: 20,
				wanted:  10,
			},
			wantScale: false,
			wantInst:  0,
			wantErr:   false,
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
			scale, inst, err := k.needScaleDown(tt.args.req, tt.args.queries, tt.args.current, tt.args.wanted)
			if (err != nil) != tt.wantErr {
				t.Errorf("needScaleToZero() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if scale != tt.wantScale {
				t.Errorf("needScaleToZero() scale = %v, wantScale %v", scale, tt.wantScale)
			}
			if tt.wantScale && inst != tt.wantInst {
				t.Errorf("needScaleToZero() inst = %v, wantInst %v", inst, tt.wantInst)
			}

		})
	}
}

func TestKubeClientAutoscaler_needScaleUp(t *testing.T) {
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
		want    int
		wantErr bool
	}{
		{
			name:   "no scale no dynamic enabled",
			fields: fields{state: MemoryState()},
			args: args{
				req: KubeRequest{
					Coordinator: testUtil.MustUrl("http://coordinator.local"),
					ScaleAfter:  10 * time.Second,
					Max:         5,
					DynamicScale: configuration.AutoscalerDynamicScale{
						Enabled: false,
					},
				},
				queries: nil,
			},
			want:    0,
			wantErr: false,
		},
		{
			name:   "scale no dynamic enabled (waiting query)",
			fields: fields{state: MemoryState()},
			args: args{
				req: KubeRequest{
					Coordinator: testUtil.MustUrl("http://coordinator.local"),
					ScaleAfter:  10 * time.Second,
					Max:         5,
					DynamicScale: configuration.AutoscalerDynamicScale{
						Enabled: false,
					},
				},
				queries: trino.QueryList{
					{
						State:       StateWaitingForResources,
						SessionUser: "test",
						QueryStats: trino.QueryStats{
							EndTime: time.Now().Add(-1 * time.Hour),
						},
					},
				},
			},
			want:    5,
			wantErr: false,
		},
		{
			name:   "no scale no dynamic enabled (running query)",
			fields: fields{state: MemoryState()},
			args: args{
				req: KubeRequest{
					Coordinator: testUtil.MustUrl("http://coordinator.local"),
					ScaleAfter:  10 * time.Second,
					Max:         5,
					DynamicScale: configuration.AutoscalerDynamicScale{
						Enabled: false,
					},
				},
				queries: trino.QueryList{
					{
						State:       StateRunning,
						SessionUser: "test",
						QueryStats: trino.QueryStats{
							EndTime: time.Now().Add(-1 * time.Hour),
						},
					},
				},
			},
			want:    0,
			wantErr: false,
		},
		{
			name:   "don't scale up when no queries and no state are present",
			fields: fields{state: MemoryState()},
			args: args{
				req: KubeRequest{
					Coordinator: testUtil.MustUrl("http://coordinator.local"),
					ScaleAfter:  10 * time.Second,
					DynamicScale: configuration.AutoscalerDynamicScale{
						Enabled: true,
						Default: 5,
						Rules: []configuration.AutoscalerDynamicScaleRule{
							{
								Regexp:    "etl-*",
								Instances: 8,
							},
						},
					},
				},
				queries: nil,
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "don't scale up when no queries but default bigger than state",
			fields: fields{state: mockState{
				getInstances: func(clusterID string) (int32, error) {
					return 0, nil
				},
			}},
			args: args{
				req: KubeRequest{
					Coordinator: testUtil.MustUrl("http://coordinator.local"),
					ScaleAfter:  10 * time.Second,
					DynamicScale: configuration.AutoscalerDynamicScale{
						Enabled: true,
						Default: 5,
						Rules: []configuration.AutoscalerDynamicScaleRule{
							{
								Regexp:    "etl-*",
								Instances: 8,
							},
						},
					},
				},
				queries: nil,
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "scale up when queries not trigger rule but default",
			fields: fields{state: mockState{
				getInstances: func(clusterID string) (int32, error) {
					return 0, nil
				},
			}},
			args: args{
				req: KubeRequest{
					Coordinator: testUtil.MustUrl("http://coordinator.local"),
					ScaleAfter:  10 * time.Second,
					DynamicScale: configuration.AutoscalerDynamicScale{
						Enabled: true,
						Default: 5,
						Rules: []configuration.AutoscalerDynamicScaleRule{
							{
								Regexp:    "etl-*",
								Instances: 8,
							},
						},
					},
				},
				queries: trino.QueryList{
					{
						State:       StateRunning,
						SessionUser: "aaaaa",
						QueryStats: trino.QueryStats{
							EndTime: time.Now().Add(-1 * time.Hour),
						},
					},
				},
			},
			want:    5,
			wantErr: false,
		},
		{
			name: "scale up when queries trigger one rule ",
			fields: fields{state: mockState{
				getInstances: func(clusterID string) (int32, error) {
					return 0, nil
				},
			}},
			args: args{
				req: KubeRequest{
					Coordinator: testUtil.MustUrl("http://coordinator.local"),
					ScaleAfter:  10 * time.Second,
					DynamicScale: configuration.AutoscalerDynamicScale{
						Enabled: true,
						Default: 5,
						Rules: []configuration.AutoscalerDynamicScaleRule{
							{
								Regexp:    "etl-*",
								Instances: 8,
							},
						},
					},
				},
				queries: trino.QueryList{
					{
						State:       StateRunning,
						SessionUser: "etl-test",
						QueryStats: trino.QueryStats{
							EndTime: time.Now().Add(-1 * time.Hour),
						},
					},
				},
			},
			want:    8,
			wantErr: false,
		},
		{
			name: "scale up when queries trigger one rule and one default get greater (rules)",
			fields: fields{state: mockState{
				getInstances: func(clusterID string) (int32, error) {
					return 0, nil
				},
			}},
			args: args{
				req: KubeRequest{
					Coordinator: testUtil.MustUrl("http://coordinator.local"),
					ScaleAfter:  10 * time.Second,
					DynamicScale: configuration.AutoscalerDynamicScale{
						Enabled: true,
						Default: 5,
						Rules: []configuration.AutoscalerDynamicScaleRule{
							{
								Regexp:    "etl-*",
								Instances: 8,
							},
						},
					},
				},
				queries: trino.QueryList{
					{
						State:       StateRunning,
						SessionUser: "etl-test",
						QueryStats: trino.QueryStats{
							EndTime: time.Now().Add(-1 * time.Hour),
						},
					},
					{
						State:       StateWaitingForResources,
						SessionUser: "test",
						QueryStats: trino.QueryStats{
							EndTime: time.Now().Add(-1 * time.Hour),
						},
					},
				},
			},
			want:    8,
			wantErr: false,
		},
		{
			name: "scale up when queries trigger one rule and one default get greater(default)",
			fields: fields{state: mockState{
				getInstances: func(clusterID string) (int32, error) {
					return 0, nil
				},
			}},
			args: args{
				req: KubeRequest{
					Coordinator: testUtil.MustUrl("http://coordinator.local"),
					ScaleAfter:  10 * time.Second,
					DynamicScale: configuration.AutoscalerDynamicScale{
						Enabled: true,
						Default: 5,
						Rules: []configuration.AutoscalerDynamicScaleRule{
							{
								Regexp:    "etl-*",
								Instances: 2,
							},
						},
					},
				},
				queries: trino.QueryList{
					{
						State:       StateRunning,
						SessionUser: "etl-test",
						QueryStats: trino.QueryStats{
							EndTime: time.Now().Add(-1 * time.Hour),
						},
					},
					{
						State:       StateWaitingForResources,
						SessionUser: "test",
						QueryStats: trino.QueryStats{
							EndTime: time.Now().Add(-1 * time.Hour),
						},
					},
				},
			},
			want:    5,
			wantErr: false,
		},
		{
			name: "scale up to default no query trigger rule",
			fields: fields{state: mockState{
				getInstances: func(clusterID string) (int32, error) {
					return 0, nil
				},
			}},
			args: args{
				req: KubeRequest{
					Coordinator: testUtil.MustUrl("http://coordinator.local"),
					ScaleAfter:  10 * time.Second,
					DynamicScale: configuration.AutoscalerDynamicScale{
						Enabled: true,
						Default: 5,
						Rules: []configuration.AutoscalerDynamicScaleRule{
							{
								Regexp:    "etl-*",
								Instances: 3,
							},
						},
					},
				},
				queries: trino.QueryList{
					{
						State:       StateRunning,
						SessionUser: "test 1",
						QueryStats: trino.QueryStats{
							EndTime: time.Now().Add(-1 * time.Hour),
						},
					},
					{
						State:       StateWaitingForResources,
						SessionUser: "test 2",
						QueryStats: trino.QueryStats{
							EndTime: time.Now().Add(-1 * time.Hour),
						},
					},
				},
			},
			want:    5,
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

			got, err := k.desiredInstances(tt.args.req, tt.args.queries)
			if (err != nil) != tt.wantErr {
				t.Errorf("desiredInstances() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("desiredInstances() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKubeClientAutoscaler_currentInstances(t *testing.T) {
	type fields struct {
		client   kubernetes.Interface
		trinoApi trino.Api
		state    State
	}
	type args struct {
		req KubeRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "state to 0 return state",
			fields: fields{state: mockState{
				getInstances: func(clusterID string) (int32, error) {
					return 0, nil
				},
			}},
			args: args{
				req: KubeRequest{
					Coordinator: testUtil.MustUrl("http://coordinator.local"),
				},
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "state err",
			fields: fields{state: mockState{
				getInstances: func(clusterID string) (int32, error) {
					return 0, fmt.Errorf("error on state")
				},
			}},
			args: args{
				req: KubeRequest{
					Coordinator: testUtil.MustUrl("http://coordinator.local"),
				},
			},
			want:    0,
			wantErr: true,
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

			got, err := k.currentInstances(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("desiredInstances() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("desiredInstances() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKubeClientAutoscaler_execute(t *testing.T) {
	req := KubeRequest{
		Coordinator: testUtil.MustUrl("http://coordinator.local"),
		Min:         0,
		Max:         5,
	}
	instance := int32(0)
	instanceLast := int32(0)

	MS := mockState{
		setInstances: func(clusterID string, i int32) error {
			instance = i
			return nil
		},
		getInstances: func(clusterID string) (int32, error) {
			return instance, nil
		},
		setTime: func(clusterID string, t time.Time) error {
			return nil
		},
		getTime: func(clusterID string) (time.Time, error) {
			return time.Now().Add(-11 * time.Minute), nil
		},
		setLastScaleUp: func(clusterID string, i int32, t time.Time) error {
			instanceLast = i
			return nil
		},
		getLastScaleUp: func(clusterID string) (int32, time.Time, error) {
			return instanceLast, time.Now().Add(-11 * time.Minute), nil
		},
	}

	k := &KubeClientAutoscaler{
		client:   nil,
		trinoApi: nil,
		state:    MS,
		logger:   logging.Noop(),
	}

	queries := trino.QueryList{
		{
			State: StateWaitingForResources,
			QueryStats: trino.QueryStats{
				EndTime: time.Unix(0, 0),
			},
		},
	}

	cInstances, err := k.currentInstances(req)
	require.NoError(t, err)
	rInstances, err := k.desiredInstances(req, queries)
	require.NoError(t, err)

	require.NotEqual(t, cInstances, rInstances)

	instance = int32(req.Max)

	cInstances, err = k.currentInstances(req)
	require.NoError(t, err)
	rInstances, err = k.desiredInstances(req, queries)
	require.NoError(t, err)

	require.Equal(t, cInstances, rInstances)

	down, err := k.needScaleToZero(req, queries)
	require.NoError(t, err)

	require.False(t, down)

	down, i, err := k.needScaleDown(req, queries, cInstances, rInstances)
	require.NoError(t, err)

	require.False(t, down)
	require.Equal(t, 0, i)

}
