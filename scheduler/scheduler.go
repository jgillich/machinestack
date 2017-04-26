package scheduler

import (
	"fmt"
	"io"

	"gitlab.com/faststack/machinestack/config"
	"gitlab.com/faststack/machinestack/driver"
)

// Scheduler is what distirbutes Driver requests
type Scheduler interface {
	Create(name, image, driverName string, attrs driver.MachineAttributes) (string, error)
	Delete(name, driverName, node string) error
	Session(name, driverName, node string, stdin io.ReadCloser, stdout io.WriteCloser, control chan driver.ControlMessage, width, height int) error
}

// NewScheduler creates a new Scheduler
func NewScheduler(name string, options *config.DriverOptions) (Scheduler, error) {
	switch name {
	case "local":
		return NewLocalScheduler(options)
	case "consul":
		return NewConsulScheduler(options)
	default:
		return nil, fmt.Errorf("unknown scheduler '%s'", name)
	}
}
