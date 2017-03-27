package handler

import (
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func (h *Handler) CreateMachine(c echo.Context) error {

	claims := c.Get("user").(*jwt.Token).Claims.(*JwtClaims)

	count, err := h.db.Model(&Machine{}).Where("machine.owner = ?", claims.Name).Count()
	if err != nil {
		return err
	}

	if count >= claims.AppMetadata.Quota.Instances {
		return c.String(http.StatusMethodNotAllowed, "quota exceeded")
	}

	machine := new(Machine)
	if err := c.Bind(machine); err != nil {
		return err
	}

	node, err := h.sched.Create(machine.Name, machine.Image, machine.Driver)
	if err != nil {
		return err
	}

	if err = h.db.Insert(&Machine{
		Name:  machine.Name,
		Image: machine.Image,
		Owner: claims.Name,
		Node:  node,
	}); err != nil {
		// TODO machine still exists here, what to do?
		return err
	}

	return c.String(http.StatusCreated, "created")
}
