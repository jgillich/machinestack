package api

import (
	"io"
	"net/http"
	"time"

	"github.com/dchest/uniuri"
	"github.com/go-pg/pg"
	"github.com/google/jsonapi"
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

// SessionCreateRequest defines the data structure of a SessionCreate request
type SessionCreateRequest struct {
	Name   string `jsonapi:"attr,name"`
	Width  int    `jsonapi:"attr,width"`
	Height int    `jsonapi:"attr,height"`
}

// SessionCreateResponse defines the data structure of a SessionCreate response
type SessionCreateResponse struct {
	ID string `jsonapi:"primary,sessions"`
}

// SessionCreate creates a new exec session
func (h *Handler) SessionCreate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	claims, err := jwtmiddleware.ContextClaims(r)
	if err != nil {
		WriteOneError(w, http.StatusUnauthorized, UnauthorizedError)
		return
	}

	create := new(SessionCreateRequest)
	if err := jsonapi.UnmarshalPayload(r.Body, create); err != nil {
		WriteOneError(w, http.StatusBadRequest, BadRequestError)
		return
	}

	var machine model.Machine
	if err := h.DB.Model(&machine).Where("name = ?", create.Name).Select(); err != nil {
		if err != pg.ErrNoRows {
			WriteInternalError(w, "session create: db error", err)
			return
		}
		WriteOneError(w, http.StatusNotFound, ResourceNotFoundError)
		return
	}

	if machine.UserID != int64(claims["id"].(float64)) {
		WriteOneError(w, http.StatusForbidden, AccessDeniedError)
		return
	}

	inr, inw := io.Pipe()
	outr, outw := io.Pipe()
	control := make(chan driver.ControlMessage)

	if err := h.Scheduler.Session(machine.Name, machine.Driver, machine.Node, inr, outw, control, create.Width, create.Height); err != nil {
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

	WriteOne(w, http.StatusCreated, &SessionCreateResponse{ID: id})
	return
}
