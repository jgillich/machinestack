package scheduler

import (
	"errors"
	"io"

	"github.com/dchest/uniuri"

	"fmt"

	"gitlab.com/faststack/machinestack/config"
	"gitlab.com/faststack/machinestack/driver"
)

// MockScheduler stores machines in memory and does not call a driver
type MockScheduler struct {
	Machines map[string]string
}

// NewMockScheduler creates a new MockScheduler
func NewMockScheduler(options *config.DriverOptions) (Scheduler, error) {
	return &MockScheduler{
		Machines: make(map[string]string),
	}, nil
}

// Create creates a new machine
func (s *MockScheduler) Create(name, image, driverName string, attrs driver.MachineAttributes) (string, error) {
	if _, ok := s.Machines[name]; ok {
		return "", fmt.Errorf("duplicate machine '%s'", name)
	}

	s.Machines[name] = uniuri.New()

	return s.Machines[name], nil
}

// Delete deletes a machine
func (s *MockScheduler) Delete(name, driverName, node string) error {

	if id, ok := s.Machines[name]; !ok {
		return fmt.Errorf("missing machine '%s'", name)
	} else if id != node {
		return errors.New("node mismatch")
	}

	delete(s.Machines, name)

	return nil
}

// Session creates an new exec session
func (s *MockScheduler) Session(name, driverName, node string, stdin io.ReadCloser, stdout io.WriteCloser, control chan driver.ControlMessage, width, height int) error {
	return errors.New("not implemented")
}
