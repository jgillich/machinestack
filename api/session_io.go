package api

import (
	"io"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

// SessionIO exposes the io socket
func (h *Handler) SessionIO(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	e, ok := execs[id]
	if !ok {
		WriteOneError(w, http.StatusNotFound, ResourceNotFoundError)
		return
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// TODO cors
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		WriteInternalError(w, "session io: error upgrading connection", err)
		return
	}
	defer conn.Close()

	go writePump(e.r, conn)
	readPump(e.w, conn)
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
