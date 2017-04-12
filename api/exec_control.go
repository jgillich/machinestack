package api

import (
	"net/http"

	"github.com/faststackco/machinestack/driver"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

// ExecControl exposes the control socket
func (h *Handler) ExecControl(c echo.Context) error {
	id := c.Param("id")
	exec, ok := execs[id]
	if !ok {
		return Error(c, http.StatusNotFound, "exec '%s' not found", id)
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// TODO
			return true
		},
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return nil
	}
	defer conn.Close()

	for {
		var msg driver.ControlMessage
		if err := conn.ReadJSON(&msg); err != nil {
			return nil
		}

		exec.control <- msg
	}
}
