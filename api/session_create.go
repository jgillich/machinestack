package api

import (
	"io"
	"net/http"
	"time"

	"github.com/dchest/uniuri"
	"github.com/go-pg/pg"
	jwtmiddleware "github.com/jgillich/jwt-middleware"
	"github.com/julienschmidt/httprouter"
	"gitlab.com/faststack/machinestack/driver"
	"gitlab.com/faststack/machinestack/model"
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

// SessionCreateResponse defines the data structure of a ExecCreate response
type SessionCreateResponse struct {
	ID string `jsonapi:"primary,sessions"`
}

// SessionCreate creates a new exec session
func (h *Handler) SessionCreate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	name := params.ByName("name")

	claims, err := jwtmiddleware.ContextClaims(r)
	if err != nil {
		WriteOneError(w, http.StatusUnauthorized, UnauthorizedError)
		return
	}

	var machine model.Machine
	if err := h.DB.Model(&machine).Where("name = ?", name).Select(); err != nil {
		if err != pg.ErrNoRows {
			WriteInternalError(w, "session create: db error", err)
			return
		}
		WriteOneError(w, http.StatusNotFound, ResourceNotFoundError)
		return
	}

	if machine.UserID != claims["name"] {
		WriteOneError(w, http.StatusUnauthorized, AccessDeniedError)
		return
	}

	inr, inw := io.Pipe()
	outr, outw := io.Pipe()
	control := make(chan driver.ControlMessage)

	if err := h.Scheduler.Exec(machine.Name, machine.Driver, machine.Node, inr, outw, control); err != nil {
		WriteInternalError(w, "session create: scheduler exec eror", err)
		return
	}

	id := uniuri.New()

	execs[id] = exec{
		w:       inw,
		r:       outr,
		control: control,
		created: time.Now(),
	}

	WriteOne(w, http.StatusCreated, SessionCreateResponse{ID: id})
	return
}
