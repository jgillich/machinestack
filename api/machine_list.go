package api

import (
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"

	"gitlab.com/faststack/machinestack/model"
)

// MachineList lists all machines of a user
func (h *Handler) MachineList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	claims := r.Context().Value("user").(jwt.Token).Claims.(jwt.MapClaims)

	machines := []model.Machine{}
	if err := h.DB.Model(&machines).Where("user_id = ?", claims["id"]).Select(); err != nil {
		WriteInternalError(w, "session info: db error", err)
		return
	}

	WriteOne(w, http.StatusOK, machines)
}
