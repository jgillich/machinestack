package handler

import (
	"net/http"

	"github.com/labstack/echo"
)

// MachineList lists all machines of a user
func (h *Handler) MachineList(c echo.Context) error {

	claims := getJwtClaims(c)

	var machines []Machine
	if err := h.db.Model(&machines).Where("owner = ?", claims.Name).Select(); err != nil {
		return err
	}

	return Data(c, http.StatusOK, machines)
}
