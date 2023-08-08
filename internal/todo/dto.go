package todo

import "github.com/timdevlet/todo/internal/helpers"

type TodoDto struct {
	UUID      string  `json:"uuid" validate:"uuid4"`
	Title     string  `json:"title" validate:"min=1,max=200"`
	IsDone    bool    `json:"done" validate:"omitempty"`
	CreatedAt string  `json:"created_at" validate:"min=1,max=200"`
	DoneAt    *string `json:"done_at" validate:"omitempty"`
}

func (t *Todo) ToDto() *TodoDto {
	return &TodoDto{
		UUID:      t.UUID,
		Title:     t.Title,
		IsDone:    t.IsDone(),
		CreatedAt: helpers.Default(t.CreatedAt, ""),
	}
}
