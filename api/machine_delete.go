package api

import (
	"net/http"

	"github.com/go-pg/pg"
	"github.com/labstack/echo"
	"gitlab.com/faststack/machinestack/model"
)

// MachineDelete deletes a machine
func (h *Handler) MachineDelete(c echo.Context) error {

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

	if err := h.sched.Delete(name, machine.Driver, machine.Node); err != nil {
		return err
	}

	if err := h.db.Delete(&machine); err != nil {
		return err
	}

	return Message(c, http.StatusOK, "deleted")
}
