package scheduler

import (
	"testing"

	"github.com/faststackco/machinestack/config"
	"github.com/faststackco/machinestack/driver"
	"github.com/stretchr/testify/assert"
)

func TestLocalScheduler(t *testing.T) {
	options := config.DriverOptions{
		"lxd.remote": "unix://",
	}

	sched, err := NewLocalScheduler(&options)
	assert.NoError(t, err)

	node, err := sched.Create("TestLocalScheduler", "ubuntu/trusty", "lxd", driver.MachineAttributes{CPU: 1, RAM: 1})
	assert.NoError(t, err)

	assert.NoError(t, sched.Delete("TestLocalScheduler", "lxd", node))
}
