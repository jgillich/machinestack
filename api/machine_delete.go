package api

import (
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-pg/pg"
	"github.com/julienschmidt/httprouter"
	"gitlab.com/faststack/machinestack/model"
)

// MachineDelete deletes a machine
func (h *Handler) MachineDelete(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	name := params.ByName("name")
	claims := r.Context().Value("user").(jwt.Token).Claims.(jwt.MapClaims)

	var machine model.Machine
	if err := h.DB.Model(&machine).Where("name = ?", name).Select(); err != nil {
		if err != pg.ErrNoRows {
			WriteInternalError(w, "session delete: db error", err)
			return
		}
		WriteOneError(w, http.StatusNotFound, ResourceNotFoundError)
		return
	}

	if machine.User != claims["id"] {
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
