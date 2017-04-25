package api

import (
	"net/http"

	"github.com/go-pg/pg"
	jwtmiddleware "github.com/jgillich/jwt-middleware"
	"github.com/julienschmidt/httprouter"
	"gitlab.com/faststack/machinestack/model"
)

// MachineInfo return info about a machine
func (h *Handler) MachineInfo(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	name := params.ByName("name")

	claims, err := jwtmiddleware.ContextClaims(r)
	if err != nil {
		WriteOneError(w, http.StatusUnauthorized, UnauthorizedError)
		return
	}

	var machine model.Machine
	if err := h.DB.Model(&machine).Where("name = ?", name).Select(); err != nil {
		if err != pg.ErrNoRows {
			WriteInternalError(w, "session info: db error", err)
			return
		}
		WriteOneError(w, http.StatusNotFound, ResourceNotFoundError)
		return
	}

	if machine.UserID != claims["id"] {
		WriteOneError(w, http.StatusUnauthorized, AccessDeniedError)
		return
	}

	WriteOne(w, http.StatusOK, &machine)
}
