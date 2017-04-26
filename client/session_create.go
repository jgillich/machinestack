package client

import (
	"bytes"
	"encoding/json"
	"fmt"

	"gitlab.com/faststack/machinestack/api"

	"github.com/google/jsonapi"
)

// SessionCreate creates a new exec session
func (c *Client) SessionCreate(machineName string, width, height int) (string, error) {
	create := api.SessionCreateRequest{Name: machineName, Width: width, Height: height}

	payload, err := jsonapi.MarshalOne(&create)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(payload)

	r, err := c.request("POST", "/session", buf)
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
