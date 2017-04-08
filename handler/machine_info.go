package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

// MachineInfo return info about a machine
func (h *Handler) MachineInfo(c echo.Context) error {

	name := c.Param("name")
	claims := getJwtClaims(c)

	var machine Machine
	if err := h.db.Model(&machine).Where("machine.name = ?", name).Select(); err != nil {
		return err
	}

	if machine.User != claims.Name {
		return c.String(http.StatusBadRequest, fmt.Sprintf("machine '%s' is not owned by '%s'", name, claims.Name))
	}

	return c.JSON(http.StatusOK, machine)
}
