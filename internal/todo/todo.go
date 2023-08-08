package todo

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

type TodoService struct {
	repo ITodoRepository

	todoOwner map[string]string
}

func NewTodoService(repo ITodoRepository) ITodo {
	return &TodoService{
		repo: repo,
		// @todo: not save for go routines
		todoOwner: make(map[string]string),
	}
}

type ITodo interface {
	InsertTodo(t Todo) (string, error)
	GetUserTodos(ownerUuid string, filter TodoFilter) ([]Todo, error)
	UpdateTitle(uuid string, title string) error
	DoneTodo(uuid string, done bool) error
	DeleteTodo(uuid string) error
	FetchByUuid(uuid string) (*Todo, error)
	AllowPatchTodo(uuid string, ownerUuid string) (bool, error)
}

// ----------------------------

func (ts *TodoService) InsertTodo(t Todo) (string, error) {
	return ts.repo.InsertTodo(t)
}

// ----------------------------

func (t *TodoService) GetUserTodos(ownerUuid string, filter TodoFilter) ([]Todo, error) {
	if filter.Limit == 0 {
		filter.Limit = 1000
	}

	if filter.OwnerUuid == "" {
		filter.OwnerUuid = ownerUuid
	}

	items, err := t.repo.FetchTodos(filter)
	if err != nil {
		log.Error(err)
		return []Todo{}, err
	}

	return items, nil
}

// ----------------------------

func (t *TodoService) FetchByUuid(uuid string) (*Todo, error) {
	return t.repo.FetchByUuid(uuid)
}

// ----------------------------

func (t *TodoService) AllowPatchTodo(uuid string, ownerUuid string) (bool, error) {
	// this is crazy code.
	// add cache invalidation and thread save! Rewrite to sync.Map or LRU
	// cachedOwnerUuid, found := t.todoOwner[uuid]
	// if found {
	// 	if cachedOwnerUuid != ownerUuid {
	// 		return false, errors.New("cant putch todo of another user")
	// 	}
	// }

	//

	todo, err := t.repo.FetchByUuid(uuid)
	if err != nil {
		return false, err
	}

	// cache
	t.todoOwner[uuid] = todo.OwnerUuid

	if todo.OwnerUuid != ownerUuid {
		return false, errors.New("cant putch todo of another user")
	}

	if todo.DeletedAt != nil {
		return false, errors.New("cant putch deleted todo")
	}

	return true, nil
}

// ----------------------------

func (t *TodoService) UpdateTitle(uuid string, title string) error {
	mp := make(map[string]any)
	mp["title"] = title

	err := t.repo.UpdateTodo(uuid, mp)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

// ----------------------------

func (t *TodoService) DoneTodo(uuid string, done bool) error {
	mp := make(map[string]any)

	if done {
		mp["done_at"] = "now()"
	} else {
		mp["done_at"] = nil
	}

	err := t.repo.UpdateTodo(uuid, mp)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

// ----------------------------

func (t *TodoService) DeleteTodo(uuid string) error {
	mp := make(map[string]any)
	mp["deleted_at"] = "now()"

	err := t.repo.UpdateTodo(uuid, mp)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
