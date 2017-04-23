package client

import (
	"fmt"

	"github.com/google/jsonapi"
	"gitlab.com/faststack/machinestack/model"
)

func (c *Client) MachineInfo(name string) (*model.Machine, error) {
	r, err := c.request("GET", "/machines/"+name, nil)
	if err != nil {
		return nil, err
	}

	switch r.StatusCode {
	case 200:
		var machine *model.Machine
		if err := jsonapi.UnmarshalPayload(r.Body, machine); err != nil {
			return nil, err
		}
		return machine, nil
	default:
		// TODO unmarshal error
		return nil, fmt.Errorf("unexpected response: %s", r.Status)
	}
}
