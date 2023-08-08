package todo

import (
	"errors"

	"github.com/Masterminds/squirrel"
	log "github.com/sirupsen/logrus"
	"github.com/timdevlet/todo/internal/helpers"
	"github.com/timdevlet/todo/pkg/postgres"
)

type ITodoRepository interface {
	InsertTodo(Todo) (string, error)
	FetchTodos(TodoFilter) ([]Todo, error)
	UpdateTodo(string, map[string]any) error
	FetchByUuid(string) (*Todo, error)
}

type TodoRepository struct {
	postgres *postgres.PDB
}

func NewTodoRepository(db *postgres.PDB) ITodoRepository {
	return &TodoRepository{
		postgres: db,
	}
}

// ----------------------------

func (repo *TodoRepository) FetchByUuid(uuid string) (*Todo, error) {
	statement := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("uuid, title, owner_uuid, done_at, deleted_at").
		Where(squirrel.Eq{"uuid": uuid}).
		From("todos")

	sql, values, _ := statement.ToSql()
	rows, err := repo.postgres.DB.Query(sql, values...)
	if err != nil {
		return &Todo{}, err
	}
	defer rows.Close()

	result := make([]Todo, 0)
	for rows.Next() {
		var i Todo
		err := rows.Scan(&i.UUID, &i.Title, &i.OwnerUuid, &i.DoneAt, &i.DeletedAt)
		if err != nil {
			return &Todo{}, err
		}
		result = append(result, i)
	}

	if err = rows.Err(); err != nil {
		return &Todo{}, err
	}

	if len(result) < 1 {
		return &Todo{}, errors.New("not found by uuid")
	}

	return &result[0], nil
}

func (repo *TodoRepository) UpdateTodo(uuid string, mp map[string]any) error {
	statement := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Update("todos").
		Where(squirrel.Eq{"uuid": uuid}).
		Set("updated_at", "now()").
		Suffix("RETURNING \"uuid\"")

	for field := range mp {
		statement = statement.Set(field, mp[field])
	}

	sql, values, _ := statement.ToSql()

	log.WithField("s", values).Info(sql)

	err := repo.postgres.DB.QueryRow(sql, values...).Err()

	return err
}

// ----------------------------

func (repo *TodoRepository) InsertTodo(r Todo) (string, error) {
	// @todo move to sql file
	sqlStatement := `
		INSERT INTO todos (owner_uuid, title, description)
		VALUES ($1, $2, $3)
		RETURNING uuid
	`

	// Create random todos
	if r.Title == "$random" {
		r.Title = helpers.FakeSentence(25)
	}

	uuid := ""
	err := repo.postgres.DB.QueryRow(sqlStatement, r.OwnerUuid, r.Title, helpers.FakeSentence(1000)).Scan(&uuid)

	log.
		WithField("uuid", uuid).
		WithField("t", r.Title).
		Debug("New todo")

	return uuid, err
}

// ----------------------------

func (repo *TodoRepository) FetchTodos(filter TodoFilter) ([]Todo, error) {
	statement := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("uuid, title, owner_uuid, done_at, deleted_at, created_at, updated_at").
		From("todos").
		Where(squirrel.Eq{"owner_uuid": filter.OwnerUuid}).
		Limit(uint64(filter.Limit))

	if filter.Done != nil {
		if *filter.Done {
			statement = statement.Where(squirrel.NotEq{"done_at": nil})
		} else {
			statement = statement.Where(squirrel.Eq{"done_at": nil})
		}
	}

	if filter.Uuid != nil {
		statement = statement.Where(squirrel.Eq{"uuid": filter.Uuid})
	}

	sql, values, _ := statement.ToSql()
	log.WithField("s", values).Info(sql)
	rows, err := repo.postgres.DB.Query(sql, values...)
	if err != nil {
		return []Todo{}, err
	}
	defer rows.Close()

	result := make([]Todo, 0)
	for rows.Next() {
		var i Todo
		err := rows.Scan(&i.UUID, &i.Title, &i.OwnerUuid, &i.DoneAt, &i.DeletedAt, &i.CreatedAt, &i.UpdatedAt)
		if err != nil {
			return []Todo{}, err
		}
		result = append(result, i)
	}

	if err = rows.Err(); err != nil {
		return []Todo{}, err
	}

	log.
		WithField("total", len(result)).
		WithField("owner", filter.OwnerUuid).
		Debug("FetchTodos")

	return result, nil
}
