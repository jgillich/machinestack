package handler

import (
	"fmt"
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

	fmt.Println(machines)

	return c.JSON(http.StatusOK, machines)
}
