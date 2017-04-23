package client

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/google/jsonapi"
	"gitlab.com/faststack/machinestack/model"
)

func (c *Client) MachineCreate(machine *model.Machine) error {
	payload, err := jsonapi.MarshalOne(machine)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(payload)

	r, err := c.request("POST", "/machines", buf)
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
