package handler

import (
	"net/http"

	"github.com/faststackco/machinestack/model"
	"github.com/go-pg/pg"
	"github.com/labstack/echo"
)

// MachineInfo return info about a machine
func (h *Handler) MachineInfo(c echo.Context) error {

	name := c.Param("name")
	claims := getJwtClaims(c)

	var machine model.Machine
	if err := h.db.Model(&machine).Where("name = ?", name).Select(); err != nil {
		if err != pg.ErrNoRows {
			return err
		}
		return Error(c, http.StatusNotFound, "machine '%s' was not found", name)
	}

	if machine.Owner != claims.Name {
		return Error(c, http.StatusBadRequest, "machine '%s' is not owned by '%s'", name, claims.Name)
	}

	return Data(c, http.StatusOK, machine)
}
