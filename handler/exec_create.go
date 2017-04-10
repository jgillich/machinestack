package handler

import (
	"io"
	"net/http"
	"time"

	"github.com/faststackco/machinestack/driver"
	"github.com/go-pg/pg"
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
	if err := h.db.Model(&machine).Where("name = ?", name).Select(); err != nil {
		if err != pg.ErrNoRows {
			return err
		}
		return Error(c, http.StatusNotFound, "machine '%s' was not found", name)
	}

	if machine.Owner != claims.Name {
		return Error(c, http.StatusBadRequest, "machine '%s' is not owned by '%s'", name, claims.Name)
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
