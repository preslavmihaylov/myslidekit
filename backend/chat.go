package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Message stores messages received in a chat room
type Message struct {
	Nickname string `json:"nickname"`
	Payload  string `json:"payload"`
}

var clients = []*websocket.Conn{}
var messages = []*Message{}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	clients = append(clients, conn)
	sendMsgHistory(conn, messages)

	for {
		_, payload, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		var msg Message
		err = json.Unmarshal(payload, &msg)
		if err != nil {
			log.Println("failed to unmarshal json message")
		}

		messages = append(messages, &msg)
		broadcastMsg(clients, &msg)
	}
}

func sendMsgHistory(conn *websocket.Conn, messages []*Message) {
	for _, m := range messages {
		payload, err := json.Marshal(m)
		if err != nil {
			log.Println(err)
			return
		}

		err = conn.WriteMessage(websocket.TextMessage, payload)
		if err != nil {
			log.Println(err)
		}
	}
}

func broadcastMsg(clients []*websocket.Conn, msg *Message) {
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
		return
	}

	for _, c := range clients {
		err = c.WriteMessage(websocket.TextMessage, payload)
		if err != nil {
			log.Println(err)
		}
	}
}
