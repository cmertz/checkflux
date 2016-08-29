package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	ws "github.com/gorilla/websocket"
)

const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

type websocket struct {
	result   <-chan string
	upgrader ws.Upgrader
	ping     <-chan time.Time
	conn     *ws.Conn
}

func (s websocket) handleRead() {
	defer s.conn.Close()
	s.conn.SetReadLimit(512)
	s.conn.SetReadDeadline(time.Now().Add(pongWait))
	s.conn.SetPongHandler(func(string) error { s.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, _, err := s.conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (s websocket) handleWrite() {
	for {
		select {
		case js := <-s.result:
			s.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := s.conn.WriteMessage(ws.TextMessage, []byte(js)); err != nil {
				fmt.Println(err)
				return
			}

		case <-s.ping:
			s.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := s.conn.WriteMessage(ws.PingMessage, []byte{}); err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func (s websocket) handleWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(ws.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}

	s.conn = conn

	go s.handleWrite()
	s.handleRead()
}

// New returns a new http handler function to manage websocket connections
// It takes a channel of strings as input for the data to send over
// the actual websocket.
func wsHandler(results <-chan string) func(http.ResponseWriter, *http.Request) {
	res := websocket{}

	res.result = results

	res.upgrader = ws.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	res.ping = time.Tick(pingPeriod)

	return res.handleWebsocket
}

// New returns a json result channel
func newResultChan(results <-chan result) <-chan string {
	js := make(chan string)

	go func() {
		for {
			select {
			case res := <-results:
				r := map[string]interface{}{
					"id":    res.ID,
					"state": stateLabels[res.state],
				}

				// ignore error
				json, _ := json.Marshal(r)

				js <- string(json)
			}
		}
	}()

	return js
}
