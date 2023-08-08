package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/timdevlet/todo/internal/configs"
	"github.com/timdevlet/todo/internal/todo"
)

func TestWeb_Todos(t *testing.T) {
	w := NewWeb(&configs.Configs{
		DB_HOST:     "postgres",
		DB_PORT:     5432,
		DB_USER:     "todo",
		DB_PASSWORD: "secret",
		DB_NAME:     "todos",
	})

	e := w.Init()

	userUuid := "1d6785e2-ed20-4a08-955d-45bb1cf0c7a3"

	t.Run("should get todo", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/todos/", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		c.Set("userUuid", userUuid)

		// Assertions
		if assert.NoError(t, w.TodosGet(c)) {

			var collection CollectionDto[todo.Todo]

			json.Unmarshal(rec.Body.Bytes(), &collection)
			assert.GreaterOrEqual(t, collection.Total, 0)
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	todoUuids := []string{}

	t.Run("should store todo", func(t *testing.T) {
		payload := `{"title":"From test"}`

		//

		req := httptest.NewRequest(http.MethodPost, "/todos/", strings.NewReader(payload))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		c.Set("userUuid", userUuid)

		// Assertions
		if assert.NoError(t, w.TodosStore(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var response SuccessInsertResponse

			json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NotEmpty(t, response.Uuid, "response uuid is empty")

			todoUuids = append(todoUuids, response.Uuid)
		}
	})

	t.Run("should done todo", func(t *testing.T) {
		payload := `{"uuid":"` + todoUuids[0] + `"}`
		url := "/todos/done/"

		//

		req := httptest.NewRequest(http.MethodPost, url, strings.NewReader(payload))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		c.Set("userUuid", userUuid)
		c.Set("command", "done")

		// Assertions
		if assert.NoError(t, w.TodosPatch(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
}
