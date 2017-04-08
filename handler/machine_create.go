package handler

import (
	"net/http"

	"github.com/faststackco/machinestack/driver"
	"github.com/labstack/echo"
)

// MachineCreate creates a new machine
func (h *Handler) MachineCreate(c echo.Context) error {

	claims := getJwtClaims(c)

	count, err := h.db.Model(&Machine{}).Where("machine.user = ?", claims.Name).Count()
	if err != nil {
		return err
	}

	if count >= claims.MachineQuota.Instances {
		return c.String(http.StatusMethodNotAllowed, "quota exceeded")
	}

	machine := new(Machine)
	if err := c.Bind(machine); err != nil {
		return err
	}

	attrs := driver.MachineAttributes{CPU: claims.MachineQuota.CPU, RAM: claims.MachineQuota.RAM}

	node, err := h.sched.Create(machine.Name, machine.Image, machine.Driver, attrs)
	if err != nil {
		return err
	}

	if err = h.db.Insert(&Machine{
		Name:  machine.Name,
		Image: machine.Image,
		User:  claims.Name,
		Node:  node,
	}); err != nil {
		// TODO machine still exists here, what to do?
		return err
	}

	return c.String(http.StatusCreated, "created")
}
