package web

import (
	"encoding/json"
	"time"

	echo "github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/timdevlet/todo/internal/helpers"
	logs "github.com/timdevlet/todo/internal/log"
)

func LogMiddleware(w *Web) func(c echo.Context, reqBody, resBody []byte) {
	return func(c echo.Context, reqBody, resBody []byte) {
		req := c.Request()
		res := c.Response()
		start := time.Now()

		// request body to json string
		requestData := make(map[string]interface{})
		err := json.Unmarshal(reqBody, &requestData)
		if err != nil {
			log.Error(err)
		}

		id := req.Header.Get(echo.HeaderXRequestID)
		if id == "" {
			id = res.Header().Get(echo.HeaderXRequestID)
		}
		reqSize := req.Header.Get(echo.HeaderContentLength)
		if reqSize == "" {
			reqSize = "0"
		}

		w.LogService.InsertLog(logs.Log{
			ID:         id,
			IP:         c.RealIP(),
			Host:       req.Host,
			Method:     req.Method,
			RequestURI: req.RequestURI,
			Status:     res.Status,
			Agent:      req.UserAgent(),
			Referer:    req.Referer(),
			Start:      start,
			Stop:       start,
			Request:    requestData,
		})
	}
}

func SetAdminUserMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		uuid := c.Request().Header.Get("UserUuid")

		// @todo: hack delete it
		if !helpers.IsValidUUID(uuid) {
			uuid = "1d6785e2-ed20-4a08-955d-45bb1cf0c73b"
			log.Warn("Set user uuid to default")
		}

		c.Set("userUuid", uuid)

		return next(c)
	}
}

func ComandMiddleware(cmd string) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("command", cmd)

			return next(c)
		}
	}
}
