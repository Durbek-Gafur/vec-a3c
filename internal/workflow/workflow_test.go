package workflow

import (
	"context"
	"reflect"
	"testing"
	s "vec-node/internal/store"
)

func TestService_UpdateWorkflow(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int
		wf  *s.Workflow
	}
	tests := []struct {
		name    string
		s       *Service
		args    args
		want    *s.Workflow
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.UpdateWorkflow(tt.args.ctx, tt.args.id, tt.args.wf)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.UpdateWorkflow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.UpdateWorkflow() = %v, want %v", got, tt.want)
			}
		})
	}
}
