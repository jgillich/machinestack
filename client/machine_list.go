package client

import (
	"fmt"
	"reflect"

	"gitlab.com/faststack/machinestack/model"

	"github.com/google/jsonapi"
)

func (c *Client) MachineList() ([]*model.Machine, error) {
	r, err := c.request("GET", "/machines", nil)
	if err != nil {
		return nil, err
	}

	switch r.StatusCode {
	case 200:
		machines, err := jsonapi.UnmarshalManyPayload(r.Body, reflect.TypeOf(new(model.Machine)))
		if err != nil {
			return nil, err
		}

		result := make([]*model.Machine, len(machines))
		for i, m := range machines {
			result[i] = m.(*model.Machine)
		}
		return result, nil
	default:
		// TODO unmarshal error
		return nil, fmt.Errorf("unexpected response: %s", r.Status)
	}
}
