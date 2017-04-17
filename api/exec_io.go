package api

import (
	"io"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

// ExecIO exposes stdin and stdout
func (h *Handler) ExecIO(c echo.Context) error {

	id := c.Param("id")
	e, ok := execs[id]
	if !ok {
		return Error(c, http.StatusNotFound, "exec '%s' not found", id)
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// TODO cors
			return true
		},
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	go writePump(e.r, conn)
	readPump(e.w, conn)

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
