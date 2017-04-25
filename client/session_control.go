package client

import (
	"strings"

	"github.com/gorilla/websocket"
	"gitlab.com/faststack/machinestack/driver"
)

// SessionControl returns a control message channel for a session
func (c *Client) SessionControl(sessionID string) (chan driver.ControlMessage, error) {
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(strings.Replace(c.url, "http", "ws", 1)+"/session/"+sessionID+"/control", nil)
	if err != nil {
		return nil, err
	}

	channel := make(chan driver.ControlMessage)

	go func() {
		defer conn.Close()
		for msg := range channel {
			if err := conn.WriteJSON(msg); err != nil {
				return
			}
		}
	}()

	return channel, nil
}
