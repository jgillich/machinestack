package handler

import (
	"fmt"
	"io"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/faststackco/machinestack/driver"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

var (
	execs map[string]exec
)

type exec struct {
	stdin   io.WriteCloser
	stdout  io.ReadCloser
	control chan driver.ControlMessage
	created time.Time
}

// CreateExec creates a new exec session
func (h *Handler) CreateExec(c echo.Context) error {

	name := c.Param("name")
	claims := c.Get("user").(*jwt.Token).Claims.(*JwtClaims)

	var machine Machine
	if err := h.db.Model(&machine).Where("machine.name = ?", name).Select(); err != nil {
		return err
	}

	if machine.Owner != claims.Name {
		return c.String(http.StatusBadRequest, fmt.Sprintf("machine '%s' is not owned by '%s'", name, claims.Name))
	}

	inr, inw := io.Pipe()
	outr, outw := io.Pipe()
	control := make(chan driver.ControlMessage)

	h.sched.Exec(machine.Name, machine.Driver, machine.Node, inr, outw, control)

	id := uuid.New().String()

	execs[id] = exec{
		stdin:   inw,
		stdout:  outr,
		control: control,
		created: time.Now(),
	}

	return c.String(http.StatusCreated, id)
}

// ExecIO exposes stdin and stdout
func (h *Handler) ExecIO(c echo.Context) error {

	id := c.Param("id")
	exec, ok := execs[id]
	if !ok {
		return c.String(http.StatusNotFound, fmt.Sprintf("exec '%s' not found", id))
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

// ExecControl exposes the control socket
func (h *Handler) ExecControl(c echo.Context) error {
	id := c.Param("id")
	exec, ok := execs[id]
	if !ok {
		return c.String(http.StatusNotFound, fmt.Sprintf("exec '%s' not found", id))
	}

	upgrader := websocket.Upgrader{}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return nil
	}
	defer conn.Close()

	for {
		var msg driver.ControlMessage
		if err := conn.ReadJSON(msg); err != nil {
			return nil
		}

		exec.control <- msg
	}
}
