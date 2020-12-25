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
	type Customer struct {
		Index     string        `json:"i"`
		Name      string        `json:"name"`
		Duration  time.Duration `json:"duration"`
		Countdown time.Duration `json:"countdown"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("err:", err)
			return
		}

		for {
			// prevent instant refresh
			ticker := time.NewTicker(50 * time.Millisecond)
			for range ticker.C {
				var cs []Customer
				for _, c := range app.customers.All() {
					cs = append(cs, Customer{
						Index:     c.ID(),
						Name:      c.Name,
						Duration:  c.Duration,
						Countdown: c.RemainingTime(),
					})
				}

				msg, err := json.Marshal(cs)
				if err != nil {
					log.Println("err:", err)
					return
				}

				if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
					log.Println(err)
					return
				}
			}
		}
	}
}
