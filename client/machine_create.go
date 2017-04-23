package client

import (
	"fmt"

	"io"

	"github.com/google/jsonapi"
	"gitlab.com/faststack/machinestack/model"
)

func (c *Client) MachineCreate(machine *model.Machine) error {
	reader, writer := io.Pipe()
	if err := jsonapi.MarshalOnePayload(writer, machine); err != nil {
		return err
	}
	r, err := c.request("POST", "/machines", reader)
	if err != nil {
		return err
	}

	switch r.StatusCode {
	case 201:
		return nil
	default:
		// TODO unmarshal error
		return fmt.Errorf("unexpected response: %s", r.Status)
	}
}
