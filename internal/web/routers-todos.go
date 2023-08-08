package web

import (
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"

	echo "github.com/labstack/echo/v4"
	"github.com/timdevlet/todo/internal/helpers"
	"github.com/timdevlet/todo/internal/todo"
)

func initTodosRoutes(a *Web, e *echo.Echo) {
	todos := e.Group("/todos")

	todos.GET("", a.TodosGet)
	todos.POST("", a.TodosStore)
	todos.PATCH("", a.TodosPatch, ComandMiddleware("title"))
	todos.PATCH("/done", a.TodosPatch, ComandMiddleware("done"))
	todos.PATCH("/undone", a.TodosPatch, ComandMiddleware("undone"))
	todos.PATCH("/delete", a.TodosPatch, ComandMiddleware("delete"))
	todos.PATCH("/title", a.TodosPatch, ComandMiddleware("title"))
}

func (a *Web) TodosGet(c echo.Context) error {
	type TodoGetRequest struct {
		Done  *bool   `json:"done"`
		Limit *int    `json:"limit" validate:"omitempty,min=1,max=100"`
		Uuid  *string `json:"uuid" validate:"omitempty,min=36,max=36"`
	}

	userId, ok := c.Get("userUuid").(string)
	if !ok {
		return errors.New("userUuid is not string")
	}

	var u TodoGetRequest
	if err := c.Bind(&u); err != nil {
		log.Error(err)
		return err
	}

	if err := c.Validate(u); err != nil {
		return err
	}

	filter := todo.TodoFilter{
		Done:      u.Done,
		Limit:     helpers.Default(u.Limit, 100),
		OwnerUuid: userId,
		Uuid:      u.Uuid,
	}

	items, err := a.TodoService.GetUserTodos(userId, filter)
	if err != nil {
		c.Error(err)
		return nil
	}

	// to Dto
	itemsDto := helpers.Map(items, func(e todo.Todo, _ int) todo.TodoDto {
		log.Debug(e.CreatedAt)
		return todo.TodoDto{
			UUID:      e.UUID,
			Title:     e.Title,
			IsDone:    e.IsDone(),
			CreatedAt: helpers.Default(e.CreatedAt, ""),
			DoneAt:    e.DoneAt,
		}
	})

	// to CollectionDto
	collection := NewCollectionDto(itemsDto, 0, filter.Limit)

	return c.JSON(http.StatusOK, collection)
}

func (a *Web) TodosStore(c echo.Context) error {
	type TodoStoreRequest struct {
		Title string `json:"title" validate:"required,min=1,max=200"`
	}

	userId := c.Get("userUuid").(string)

	var u TodoStoreRequest

	if err := c.Bind(&u); err != nil {
		log.Error(err)
		return err
	}

	if err := c.Validate(u); err != nil {
		return err
	}

	uuid, err := a.TodoService.InsertTodo(todo.Todo{
		Title:     u.Title,
		OwnerUuid: userId,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, SuccessInsertResponse{
		Uuid: uuid,
	})
}

func (a *Web) TodosPatch(c echo.Context) error {
	type TodosPatchRequest struct {
		Uuid  string `json:"uuid" validate:"required,max=36,min=36"`
		Title string `json:"title" validate:"omitempty,min=1,max=200"`
		Done  bool   `json:"done" validate:"omitempty"`
	}

	userId := c.Get("userUuid").(string)
	command := c.Get("command").(string)

	var u TodosPatchRequest

	if err := c.Bind(&u); err != nil {
		log.Error(err)
		return err
	}

	if err := c.Validate(u); err != nil {
		return err
	}

	log.
		WithField("cmd", command).
		Debug("todo patch command")

	_, err := a.TodoService.AllowPatchTodo(u.Uuid, userId)
	if err != nil {
		return echo.NewHTTPError(400, err.Error())
	}

	switch command {
	case "title":
		if err := a.TodoService.UpdateTitle(u.Uuid, u.Title); err != nil {
			return err
		}
	case "delete":
		if err := a.TodoService.DeleteTodo(u.Uuid); err != nil {
			return err
		}
	case "done":
		if err := a.TodoService.DoneTodo(u.Uuid, true); err != nil {
			return err
		}
	case "undone":
		if err := a.TodoService.DoneTodo(u.Uuid, false); err != nil {
			return err
		}

	default:
		return echo.NewHTTPError(400, "no valid command (cmd) provided")
	}

	return c.NoContent(http.StatusOK)
}
