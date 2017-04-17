package driver

import (
	"fmt"
	"io"

	"gitlab.com/faststack/machinestack/config"
	"github.com/lxc/lxd/shared/api"
)

// NewDriver creates a new driver of type name
func NewDriver(name string, options config.DriverOptions) (Driver, error) {
	switch name {
	case "lxd":
		return NewLxdDriver(options)
	default:
		return nil, fmt.Errorf("unknown driver '%s'", name)
	}
}

// Driver is what creates and executes machines
type Driver interface {
	Create(name, image string, attrs MachineAttributes) error
	Delete(name string) error
	Exec(name string, stdin io.ReadCloser, stdout io.WriteCloser, control chan ControlMessage) error
}

// ControlMessage is used to send signals like resize to machines
// TODO we probably want a generic type here in the future
type ControlMessage api.ContainerExecControl

// MachineAttributes defines custom properties used for machine creation
type MachineAttributes struct {
	RAM int
	CPU int
}
