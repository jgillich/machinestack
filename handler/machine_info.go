package handler

import (
	"net/http"

	"github.com/labstack/echo"
)

// MachineInfo return info about a machine
func (h *Handler) MachineInfo(c echo.Context) error {

	name := c.Param("name")
	claims := getJwtClaims(c)

	var machine Machine
	if err := h.db.Model(&machine).Where("name = ?", name).Select(); err != nil {
		return err
	}

	if machine.Owner != claims.Name {
		return Error(c, http.StatusBadRequest, "machine '%s' is not owned by '%s'", name, claims.Name)
	}

	return Data(c, http.StatusOK, machine)
}
