package client

import (
	"io"
	"strings"

	"github.com/gorilla/websocket"
)

// SessionIO transmits input and output for a session
func (c *Client) SessionIO(sessionID string, r io.ReadCloser, w io.WriteCloser) error {
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(strings.Replace(c.url, "http", "ws", 1)+"/session/"+sessionID+"/io", nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	go writePump(r, conn)
	readPump(w, conn)

	return nil
}

func readPump(w io.WriteCloser, conn *websocket.Conn) {
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			return
		}
		if mt != websocket.TextMessage {
			continue
		}

		if _, err := w.Write(message); err != nil {
			conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseGoingAway, "broken write pipe"))
			return
		}
	}
}

func writePump(r io.ReadCloser, conn *websocket.Conn) {
	for {
		p := make([]byte, 32*1024)

		n, err := r.Read(p)
		if err != nil {
			conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseGoingAway, "broken read pipe"))
			return
		}

		if err := conn.WriteMessage(websocket.TextMessage, p[0:n]); err != nil {
			return
		}
	}
}
