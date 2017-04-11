package api

import (
	"io"
	"net/http"
	"time"

	"github.com/faststackco/machinestack/driver"
	"github.com/faststackco/machinestack/model"
	"github.com/go-pg/pg"
	"github.com/google/uuid"
	"github.com/labstack/echo"
)

var (
	execs = make(map[string]exec)
)

type exec struct {
	w       io.WriteCloser
	r       io.ReadCloser
	control chan driver.ControlMessage
	created time.Time
}

// ExecCreateResponse defines the data structure of a ExecCreate response
type ExecCreateResponse struct {
	ID string `json:"id"`
}

// ExecCreate creates a new exec session
func (h *Handler) ExecCreate(c echo.Context) error {

	name := c.Param("name")
	claims := getJwtClaims(c)

	var machine model.Machine
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

	if err := h.sched.Exec(machine.Name, machine.Driver, machine.Node, inr, outw, control); err != nil {
		return err
	}

	id := uuid.New().String()

	execs[id] = exec{
		w:       inw,
		r:       outr,
		control: control,
		created: time.Now(),
	}

	return Data(c, http.StatusCreated, ExecCreateResponse{ID: id})
}
