package driver

import (
	"fmt"
	"io"

	"github.com/faststackco/machinestack/config"
	"github.com/lxc/lxd/shared/api"
)

type Factory func(config.DriverOptions) (Driver, error)

var BuiltinDrivers = map[string]Factory{
	"lxd": NewLxdDriver,
}

func NewDriver(name string, options config.DriverOptions) (Driver, error) {
	factory, ok := BuiltinDrivers[name]
	if !ok {
		return nil, fmt.Errorf("unknown driver '%s'", name)
	}

	return factory(options)
}

type Driver interface {
	Create(name, image string) error
	Delete(name string) error
	Exec(name string, stdin io.ReadCloser, stdout io.WriteCloser, control chan ControlMessage) error
}

// TODO we probably want a generic type here in the future
type ControlMessage api.ContainerExecControl
