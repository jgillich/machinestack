package api

import (
	"net/http"

	jwtmiddleware "github.com/jgillich/jwt-middleware"
	"github.com/julienschmidt/httprouter"

	"gitlab.com/faststack/machinestack/model"
)

// MachineList lists all machines of a user
func (h *Handler) MachineList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	claims, err := jwtmiddleware.ContextClaims(r)
	if err != nil {
		WriteOneError(w, http.StatusUnauthorized, UnauthorizedError)
		return
	}

	var machines []*model.Machine
	if err := h.DB.Model(&machines).Where("user_id = ?", claims["id"]).Select(); err != nil {
		WriteInternalError(w, "session info: db error", err)
		return
	}

	WriteMany(w, http.StatusOK, machines)
}
