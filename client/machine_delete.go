package client

import (
	"fmt"
)

func (c *Client) MachineDelete(name string) error {
	r, err := c.request("DELETE", "/machines/"+name, nil)
	if err != nil {
		return err
	}

	switch r.StatusCode {
	case 200:
		return nil
	default:
		// TODO unmarshal error
		return fmt.Errorf("unexpected response: %s", r.Status)
	}
}
