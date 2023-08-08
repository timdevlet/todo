package todo

import (
	"reflect"
	"testing"

	"github.com/timdevlet/todo/internal/helpers"
)

func TestTodo_Validate(t *testing.T) {
	type fields struct {
		UUID      string
		Title     string
		OwnerUuid string
		DoneAt    *string
		DeletedAt *string
		CreatedAt *string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
		want1  bool
	}{
		{
			name: "valid",
			fields: fields{
				UUID:      helpers.Uuid(),
				OwnerUuid: helpers.Uuid(),
				Title:     "hello",
			},
			want:  []string{},
			want1: true,
		},
		{
			name: "valid date",
			fields: fields{
				UUID:      helpers.Uuid(),
				OwnerUuid: helpers.Uuid(),
				Title:     "hello",
				DoneAt:    helpers.Ptr(helpers.DateNow()),
			},
			want:  []string{},
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Todo{
				UUID:      tt.fields.UUID,
				Title:     tt.fields.Title,
				OwnerUuid: tt.fields.OwnerUuid,
				DoneAt:    tt.fields.DoneAt,
				DeletedAt: tt.fields.DeletedAt,
				CreatedAt: tt.fields.CreatedAt,
			}
			got, got1 := tr.Validate()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Todo.Validate() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Todo.Validate() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
