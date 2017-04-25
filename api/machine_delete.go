package api

import (
	"net/http"

	"github.com/go-pg/pg"
	jwtmiddleware "github.com/jgillich/jwt-middleware"
	"github.com/julienschmidt/httprouter"
	"gitlab.com/faststack/machinestack/model"
)

// MachineDelete deletes a machine
func (h *Handler) MachineDelete(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	name := params.ByName("name")

	claims, err := jwtmiddleware.ContextClaims(r)
	if err != nil {
		WriteOneError(w, http.StatusUnauthorized, UnauthorizedError)
		return
	}

	var machine model.Machine
	if err := h.DB.Model(&machine).Where("name = ?", name).Select(); err != nil {
		if err != pg.ErrNoRows {
			WriteInternalError(w, "session delete: db error", err)
			return
		}
		WriteOneError(w, http.StatusNotFound, ResourceNotFoundError)
		return
	}

	if machine.UserID != claims["id"] {
		WriteOneError(w, http.StatusUnauthorized, AccessDeniedError)
		return
	}

	if err := h.Scheduler.Delete(name, machine.Driver, machine.Node); err != nil {
		WriteInternalError(w, "machine delete: delete failed", err)
		return
	}

	if err := h.DB.Delete(&machine); err != nil {
		WriteInternalError(w, "machine delete: db error", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
