package web

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// RemainingRealTime gives real-time remaining time to client
func RemainingRealTime(app *App) http.HandlerFunc {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("err:", err)
			return
		}

		for {
			select {
			// prevent instant refresh
			case <-time.After(50 * time.Millisecond):
				msg, err := json.Marshal(app.customers.All())
				if err != nil {
					log.Println("err:", err)
					return
				}
				app.incr.Incr()

				if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
					log.Println(err)
					return
				}
			}
		}
	}
}
