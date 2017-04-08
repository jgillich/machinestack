package handler

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/faststackco/machinestack/driver"
	"github.com/google/uuid"
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

// ExecCreate creates a new exec session
func (h *Handler) ExecCreate(c echo.Context) error {

	name := c.Param("name")
	claims := getJwtClaims(c)

	var machine Machine
	if err := h.db.Model(&machine).Where("machine.name = ?", name).Select(); err != nil {
		return err
	}

	if machine.User != claims.Name {
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
