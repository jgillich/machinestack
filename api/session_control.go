package api

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"gitlab.com/faststack/machinestack/driver"
)

// SessionControl exposes the control socket
func (h *Handler) SessionControl(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	exec, ok := execs[id]
	if !ok {
		WriteOneError(w, http.StatusNotFound, ResourceNotFoundError)
		return
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			for _, o := range h.AllowOrigins {
				if o == origin || o == "*" {
					return true
				}
			}
			return false
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		WriteInternalError(w, "session control: error upgrading connection", err)
		return
	}
	defer conn.Close()

	for {
		var msg driver.ControlMessage
		if err := conn.ReadJSON(&msg); err != nil {
			conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseGoingAway, "malformed message"))
			return
		}

		exec.control <- msg
	}
}
