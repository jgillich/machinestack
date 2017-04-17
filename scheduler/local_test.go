package scheduler

import (
	"testing"

	"gitlab.com/faststack/machinestack/config"
	"gitlab.com/faststack/machinestack/driver"
	"github.com/stretchr/testify/assert"
)

func TestLocalScheduler(t *testing.T) {
	name := "TestLocalScheduler"
	image := "ubuntu/xenial"

	options := config.DriverOptions{
		"lxd.remote": "unix://",
	}

	sched, err := NewLocalScheduler(&options)
	assert.NoError(t, err)

	node, err := sched.Create(name, image, "lxd", driver.MachineAttributes{CPU: 1, RAM: 1})
	assert.NoError(t, err)

	assert.NoError(t, sched.Delete(name, "lxd", node))
}
