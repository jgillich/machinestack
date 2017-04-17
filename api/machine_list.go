package api

import (
	"net/http"

	"gitlab.com/faststack/machinestack/model"
	"github.com/labstack/echo"
)

// MachineList lists all machines of a user
func (h *Handler) MachineList(c echo.Context) error {

	claims := getJwtClaims(c)

	machines := []model.Machine{}
	if err := h.db.Model(&machines).Where("owner = ?", claims.Name).Select(); err != nil {
		return err
	}

	return Data(c, http.StatusOK, machines)
}
