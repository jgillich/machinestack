package handler

import (
	"io"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

// ExecIO exposes stdin and stdout
func (h *Handler) ExecIO(c echo.Context) error {

	id := c.Param("id")
	exec, ok := execs[id]
	if !ok {
		return Error(c, http.StatusNotFound, "exec '%s' not found", id)
	}

	upgrader := websocket.Upgrader{}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	for {
		messageType, r, err := conn.NextReader()
		if err != nil {
			return nil
		}
		if _, err := io.Copy(exec.stdin, r); err != nil {
			return err
		}

		w, err := conn.NextWriter(messageType)
		if err != nil {
			return err
		}
		if _, err := io.Copy(w, exec.stdout); err != nil {
			return err
		}
		if err := w.Close(); err != nil {
			return err
		}
	}
}
