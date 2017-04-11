package api

import (
	"net/http"

	"github.com/faststackco/machinestack/driver"
	"github.com/faststackco/machinestack/model"
	"github.com/labstack/echo"
)

// MachineCreate creates a new machine
func (h *Handler) MachineCreate(c echo.Context) error {

	claims := getJwtClaims(c)

	if count, err := h.db.Model(&model.Machine{}).Where("owner = ?", claims.Name).Count(); err != nil {
		return err
	} else if count >= claims.MachineQuota.Instances {
		return Error(c, http.StatusForbidden, "quota exceeded")
	}

	machine := new(model.Machine)
	if err := c.Bind(machine); err != nil {
		return err
	}

	if count, err := h.db.Model(&model.Machine{}).Where("name = ?", machine.Name).Count(); err != nil {
		return err
	} else if count > 0 {
		return Error(c, http.StatusForbidden, "machine with name '%s' exists", machine.Name)
	}

	attrs := driver.MachineAttributes{CPU: claims.MachineQuota.CPU, RAM: claims.MachineQuota.RAM}

	node, err := h.sched.Create(machine.Name, machine.Image, machine.Driver, attrs)
	if err != nil {
		return err
	}

	machine.Node = node
	machine.Owner = claims.Name

	if err = h.db.Insert(&machine); err != nil {
		// TODO machine still exists here, what to do?
		return err
	}

	return Message(c, http.StatusCreated, "created")
}
