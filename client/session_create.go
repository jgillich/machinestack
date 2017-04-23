package client

import (
	"fmt"

	"gitlab.com/faststack/machinestack/api"

	"github.com/google/jsonapi"
)

func (c *Client) SessionCreate(machineName string) (string, error) {
	r, err := c.request("POST", "/machines"+machineName+"/session", nil)
	if err != nil {
		return "", err
	}

	switch r.StatusCode {
	case 201:
		var res api.SessionCreateResponse
		if err := jsonapi.UnmarshalPayload(r.Body, &res); err != nil {
			return "", err
		}
		return res.ID, nil
	default:
		// TODO unmarshal error
		return "", fmt.Errorf("unexpected response: %s", r.Status)
	}
}
