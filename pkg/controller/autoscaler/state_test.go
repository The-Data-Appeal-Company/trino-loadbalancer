package autoscaler

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestInMemory_GetClusterInstances(t *testing.T) {
	type fields struct {
		stateInstances map[string]int32
	}
	type args struct {
		clusterID string
	}
	tests := []struct {
		name                             string
		fields                           fields
		args                             args
		want                             int32
		wantErr                          bool
		wantCustomNoInstanceInStateError bool
	}{
		{
			name: "get state no error",
			fields: fields{
				stateInstances: map[string]int32{
					"test": 5,
				},
			},
			args: args{
				clusterID: "test",
			},
			want:                             5,
			wantErr:                          false,
			wantCustomNoInstanceInStateError: false,
		},
		{
			name: "get state custom error",
			fields: fields{
				stateInstances: map[string]int32{},
			},
			args: args{
				clusterID: "test",
			},
			want:                             0,
			wantErr:                          true,
			wantCustomNoInstanceInStateError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &InMemory{
				stateInstances: tt.fields.stateInstances,
			}
			got, err := i.GetClusterInstances(tt.args.clusterID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetClusterInstances() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && errors.Is(err, NoInstancesInStateError) != tt.wantCustomNoInstanceInStateError {
				t.Errorf("GetClusterInstances() error = %v, wantCustomNoInstanceInStateError %v", err, tt.wantCustomNoInstanceInStateError)
				return
			}
			if got != tt.want {
				t.Errorf("GetClusterInstances() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInMemory_LastQueryExecution(t *testing.T) {
	type fields struct {
		stateTime map[string]time.Time
	}
	type args struct {
		clusterID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name: "get state no error",
			fields: fields{
				stateTime: map[string]time.Time{
					"test": time.Unix(0, 0),
				},
			},
			args: args{
				clusterID: "test",
			},
			want:    time.Unix(0, 0),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &InMemory{
				stateTime: tt.fields.stateTime,
			}
			got, err := i.LastQueryExecution(tt.args.clusterID)
			if (err != nil) != tt.wantErr {
				t.Errorf("LastQueryExecution() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LastQueryExecution() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInMemory_SetClusterInstances(t *testing.T) {
	type fields struct {
		stateInstances map[string]int32
	}
	type args struct {
		clusterID string
		instances int32
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "set state no error",
			fields: fields{
				stateInstances: map[string]int32{},
			},
			args: args{
				clusterID: "test",
				instances: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &InMemory{
				stateInstances: tt.fields.stateInstances,
			}
			if err := i.SetClusterInstances(tt.args.clusterID, tt.args.instances); (err != nil) != tt.wantErr {
				t.Errorf("SetClusterInstances() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInMemory_SetLastQueryExecution(t *testing.T) {
	type fields struct {
		stateTime map[string]time.Time
	}
	type args struct {
		clusterID string
		t         time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "set state no error",
			fields: fields{
				stateTime: map[string]time.Time{},
			},
			args: args{
				clusterID: "test",
				t:         time.Now(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &InMemory{
				stateTime: tt.fields.stateTime,
			}
			if err := i.SetLastQueryExecution(tt.args.clusterID, tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("SetLastQueryExecution() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
