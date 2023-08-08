package todo

import "github.com/timdevlet/todo/internal/helpers"

type Todo struct {
	UUID      string  `validate:"uuid4"`
	Title     string  `validate:"min=1,max=200"`
	OwnerUuid string  `validate:"uuid4"`
	DoneAt    *string `validate:"omitempty"`
	DeletedAt *string `validate:"omitempty"`
	CreatedAt *string `validate:"omitempty"`
	UpdatedAt *string `validate:"omitempty"`
}

func (t *Todo) IsDone() bool {
	return t.DoneAt != nil
}

func (t *Todo) Validate() ([]string, bool) {
	return helpers.ValidationStruct(t)
}

// ----------------------------

type TodoFilter struct {
	Done      *bool
	OwnerUuid string
	Limit     int
	Uuid      *string
}
