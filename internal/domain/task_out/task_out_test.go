package taskin

import (
	"backoffice/internal/domain/task_out/dto"
	"context"
	"testing"
)

func Test_service_Process(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		ct      *dto.TaskOut
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var s service
			gotErr := s.Process(context.Background(), tt.ct)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Process() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Process() succeeded unexpectedly")
			}
		})
	}
}
