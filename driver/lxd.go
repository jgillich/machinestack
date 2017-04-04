package driver

import (
	"encoding/json"
	"errors"

	"io"

	"github.com/faststackco/machinestack/config"
	"github.com/gorilla/websocket"
	"github.com/lxc/lxd"
	"github.com/lxc/lxd/shared"
	"github.com/lxc/lxd/shared/api"
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
func (d *LxdDriver) Create(name, image string) error {

	profiles := []string{"default"}

	config := make(map[string]string)
	config["limits.cpu"] = "4"
	config["limits.cpu.allowance"] = "10%"
	config["limits.cpu.priority"] = "0"
	config["limits.disk.priority"] = "0"
	config["limits.memory"] = "1GB"
	config["limits.processes"] = "500"

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

// Exec creates a new exec session
func (d *LxdDriver) Exec(name string, stdin io.ReadCloser, stdout io.WriteCloser, control chan ControlMessage) error {

	controlHandlerWrapper := func(c *lxd.Client, conn *websocket.Conn) {
		for msg := range control {

			w, err := conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return // TODO log
			}

			buf, err := json.Marshal(msg)
			if err != nil {
				return // TODO log
			}
			_, err = w.Write(buf)

			w.Close()
		}
	}

	_, err := d.client.Exec(name, []string{"/bin/bash"}, nil, stdin, stdout, nil, controlHandlerWrapper, 80, 25)

	return err
}
