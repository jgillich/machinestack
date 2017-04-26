package client

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/google/jsonapi"
	"gitlab.com/faststack/machinestack/model"
)

func (c *Client) MachineCreate(machine *model.Machine) (*model.Machine, error) {
	payload, err := jsonapi.MarshalOne(machine)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(payload)

	r, err := c.request("POST", "/machines", buf)
	if err != nil {
		return nil, err
	}

	switch r.StatusCode {
	case 201:
		var res model.Machine
		if err := jsonapi.UnmarshalPayload(r.Body, &res); err != nil {
			return nil, err
		}
		return &res, nil
	default:
		// TODO unmarshal error
		return nil, fmt.Errorf("unexpected response: %s", r.Status)
	}
}
