package api

import (
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-pg/pg"
	"github.com/julienschmidt/httprouter"
	"gitlab.com/faststack/machinestack/model"
)

// MachineInfo return info about a machine
func (h *Handler) MachineInfo(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	name := params.ByName("name")
	claims := r.Context().Value("user").(jwt.Token).Claims.(jwt.MapClaims)

	var machine model.Machine
	if err := h.DB.Model(&machine).Where("name = ?", name).Select(); err != nil {
		if err != pg.ErrNoRows {
			WriteInternalError(w, "session info: db error", err)
			return
		}
		WriteOneError(w, http.StatusNotFound, ResourceNotFoundError)
		return
	}

	if machine.Owner != claims["id"] {
		WriteOneError(w, http.StatusUnauthorized, AccessDeniedError)
		return
	}

	WriteOne(w, http.StatusOK, machine)
}
