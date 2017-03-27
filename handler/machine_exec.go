package handler

import (
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

func (h *Handler) ExecMachine(c echo.Context) error {

	name := c.Param("name")
	claims := c.Get("user").(*jwt.Token).Claims.(*JwtClaims)

	var machine Machine
	if err := h.db.Model(&machine).Where("machine.name = ?", name).Select(); err != nil {
		return err
	}

	if machine.Owner != claims.Name {
		return c.String(http.StatusBadRequest, fmt.Sprintf("machine '%s' is not owned by '%s'", name, claims.Name))
	}

	h.sched.Exec()

	upgrader := websocket.Upgrader{}

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	// ...

	return nil
}
