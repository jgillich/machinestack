package api

import (
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/asaskevich/govalidator"
	jwtmiddleware "github.com/jgillich/jwt-middleware"
	"github.com/julienschmidt/httprouter"

	"github.com/google/jsonapi"
	"gitlab.com/faststack/machinestack/driver"
	"gitlab.com/faststack/machinestack/model"
)

var (
	// QuotaExceededError is returned when the machine limit for a user is reached
	QuotaExceededError = &jsonapi.ErrorObject{
		Code:   "machine_quota_exceeded",
		Title:  "Machine quota exceeded",
		Detail: "You have reached your limit of machines you are allowed to create.",
	}
	// MachineNameTakenError is returned when the machine name is taken
	MachineNameTakenError = &jsonapi.ErrorObject{
		Code:   "machine_name_taken",
		Title:  "Machine name is taken",
		Detail: "Please choose a different name.",
	}
)

// MachineCreate creates a new machine
func (h *Handler) MachineCreate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	claims, err := jwtmiddleware.ContextClaims(r)
	if err != nil {
		WriteOneError(w, http.StatusUnauthorized, UnauthorizedError)
		return
	}

	quotas := make(map[string]int)

	for _, k := range []string{"instances", "cpu", "ram"} {
		if _, ok := claims["machine_quota"]; !ok {
			WriteInternalError(w, "machine create: missing machine_quota", errors.New(""))
			return
		}

		quotas[k] = int(claims["machine_quota"].(map[string]interface{})[k].(float64))
	}

	if count, err := h.DB.Model(&model.Machine{}).Where("user_id = ?", claims["id"]).Count(); err != nil {
		WriteInternalError(w, "machine create: db error", err)
		return
	} else if count >= quotas["instances"] {
		WriteOneError(w, http.StatusForbidden, QuotaExceededError)
		return
	}

	machine := new(model.Machine)
	if err := jsonapi.UnmarshalPayload(r.Body, machine); err != nil {
		WriteOneError(w, http.StatusBadRequest, BadRequestError)
		return
	}

	if _, err := govalidator.ValidateStruct(machine); err != nil {
		e := *ValidationFailedError
		e.Detail = err.Error()
		WriteOneError(w, http.StatusBadRequest, &e)
		return
	}

	if count, err := h.DB.Model(&model.Machine{}).Where("name = ?", machine.Name).Count(); err != nil {
		WriteInternalError(w, "machine create: db error", err)
		return
	} else if count > 0 {
		WriteOneError(w, http.StatusForbidden, MachineNameTakenError)
		return
	}

	attrs := driver.MachineAttributes{CPU: quotas["cpu"], RAM: quotas["ram"]}

	node, err := h.Scheduler.Create(machine.Name, machine.Image, machine.Driver, attrs)
	if err != nil {
		WriteInternalError(w, "machine create: scheduler error", err)
		return
	}

	machine.Node = node
	machine.UserID = int64(claims["id"].(float64))

	if err = h.DB.Insert(&machine); err != nil {
		WriteInternalError(w, "machine create: db error", err)
		if err = h.Scheduler.Delete(machine.Name, machine.Driver, node); err != nil {
			logger.Error(fmt.Sprintf("machine create: cleanup of '%s' on node '%s' failed", machine.Name, node), zap.Error(err))
		}
		return
	}

	WriteOne(w, http.StatusCreated, machine)
}
