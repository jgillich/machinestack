package scheduler

import (
	"io"

	"github.com/faststackco/machinestack/config"
	"github.com/faststackco/machinestack/driver"
)

type LocalScheduler struct {
	driverOptions *config.DriverOptions
}

func NewLocalScheduler(options *config.DriverOptions) (Scheduler, error) {
	return &LocalScheduler{
		driverOptions: options,
	}, nil
}

func (c *LocalScheduler) Create(name, image, driverName string) (string, error) {
	driver, err := driver.NewDriver(name, *c.driverOptions)
	if err != nil {
		return "", err
	}

	if err := driver.Create(name, image); err != nil {
		return "", err
	}

	return "", nil
}

func (c *LocalScheduler) Delete(name, driverName, node string) error {

	driver, err := driver.NewDriver(driverName, *c.driverOptions)
	if err != nil {
		return err
	}

	if err := driver.Delete(name); err != nil {
		return err
	}

	return nil
}

func (c *LocalScheduler) Exec(name, driverName, node string, stdin io.ReadCloser, stdout io.WriteCloser, control chan driver.ControlMessage) error {
	driver, err := driver.NewDriver(driverName, *c.driverOptions)
	if err != nil {
		return err
	}

	return driver.Exec(name, stdin, stdout, control)
}
