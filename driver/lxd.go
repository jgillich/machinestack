package driver

import (
	"errors"

	"io"

	"fmt"

	"github.com/gorilla/websocket"
	"github.com/lxc/lxd"
	"github.com/lxc/lxd/shared"
	"github.com/lxc/lxd/shared/api"
	"gitlab.com/faststack/machinestack/config"
)

// LxdDriver implements the Driver interface for LXD
type LxdDriver struct {
	client *lxd.Client
}

// NewLxdDriver creates a new LxdDriver
func NewLxdDriver(options config.DriverOptions) (Driver, error) {

	remote := lxd.RemoteConfig{
		Addr:   options["lxd.remote"],
		Static: true,
		Public: false,
	}

	lxdConfig := lxd.Config{
		Remotes: map[string]lxd.RemoteConfig{
			"remote": remote,
			"images": lxd.ImagesRemote,
		},
		DefaultRemote: "remote",
	}

	client, err := lxd.NewClient(&lxdConfig, "remote")

	if err != nil {
		return nil, err
	}

	return &LxdDriver{client: client}, nil
}

// Create creates a new machine
func (d *LxdDriver) Create(name, image string, attrs MachineAttributes) error {

	profiles := []string{"default"}

	config := map[string]string{
		"limits.cpu":           string(attrs.CPU),
		"limits.cpu.allowance": "10%",
		"limits.cpu.priority":  "0",
		"limits.disk.priority": "0",
		"limits.memory":        fmt.Sprintf("%vGB", attrs.RAM),
		"limits.processes":     "500",
	}

	res, err := d.client.Init(name, "images", image, &profiles, config, nil, false)
	if err != nil {
		return err
	}

	if err := d.client.WaitForSuccess(res.Operation); err != nil {
		return err
	}

	return nil
}

// Delete deletes machine
func (d *LxdDriver) Delete(name string) error {

	container, err := d.client.ContainerInfo(name)
	if err != nil {
		return err
	}

	if container.StatusCode != 0 && container.StatusCode != api.Stopped {
		resp, err := d.client.Action(name, shared.Stop, -1, true, false)
		if err != nil {
			return err
		}

		op, err := d.client.WaitFor(resp.Operation)
		if err != nil {
			return err
		}

		if op.StatusCode == api.Failure {
			return errors.New("stopping container failed")
		}
	}

	res, err := d.client.Delete(name)
	if err != nil {
		return err
	}

	if err := d.client.WaitForSuccess(res.Operation); err != nil {
		return err
	}

	return nil
}

// Session creates a new exec session
func (d *LxdDriver) Session(name string, stdin io.ReadCloser, stdout io.WriteCloser, control chan ControlMessage, width, height int) error {

	container, err := d.client.ContainerInfo(name)
	if err != nil {
		return err
	}

	if container.StatusCode == api.Stopped {
		if resp, err := d.client.Action(name, shared.Start, -1, false, false); err != nil {
			return err
		} else if err := d.client.WaitForSuccess(resp.Operation); err != nil {
			return err
		}
	}

	controlHandlerWrapper := func(c *lxd.Client, conn *websocket.Conn) {
		for msg := range control {
			if err := conn.WriteJSON(msg); err != nil {
				return
			}
		}
	}

	go func() {
		// TODO exec unfortunately blocks, so we cannot receive the error
		_, err = d.client.Exec(name, []string{"/bin/bash"}, env, stdin, stdout, nil, controlHandlerWrapper, width, height)
		if err != nil {
			panic(err)
		}
	}()

	return nil
}
