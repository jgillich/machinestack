package handler

import (
	"net/http"

	"github.com/labstack/echo"
)

// MachineDelete deletes a machine
func (h *Handler) MachineDelete(c echo.Context) error {

	name := c.Param("name")
	claims := getJwtClaims(c)

	var machine Machine
	if err := h.db.Model(&machine).Where("name = ?", name).Select(); err != nil {
		return err
	}

	if machine.Owner != claims.Name {
		return Error(c, http.StatusBadRequest, "machine '%s' is not owned by '%s'", name, claims.Name)
	}

	if err := h.sched.Delete(name, machine.Driver, machine.Node); err != nil {
		return err
	}

	return Message(c, http.StatusOK, "deleted")
}
