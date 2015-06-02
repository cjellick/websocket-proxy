package proxy

import (
	"strings"

	"code.google.com/p/go-uuid/uuid"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"

	"github.com/rancherio/websocket-proxy/common"
)

type Multiplexer struct {
	backendId         string
	messagesToBackend chan string
	clients           map[string]chan<- string
}

func (m *Multiplexer) initializeClient() (string, <-chan string) {
	msgKey := uuid.New()
	clientChan := make(chan string)
	m.clients[msgKey] = clientChan
	return msgKey, clientChan
}

func (m *Multiplexer) connect(msgKey, url string) {
	message := common.FormatMessage(msgKey, common.Connect, url)
	m.messagesToBackend <- message
}

func (m *Multiplexer) send(msgKey, msg string) {
	message := common.FormatMessage(msgKey, common.Body, msg)
	m.messagesToBackend <- message
}

func (m *Multiplexer) sendClose(msgKey string) {
	message := common.FormatMessage(msgKey, common.Close, "")
	m.messagesToBackend <- message
}

func (m *Multiplexer) closeConnection(msgKey string) {
	m.sendClose(msgKey)
	delete(m.clients, msgKey)
}

func (m *Multiplexer) routeMessages(ws *websocket.Conn) {
	// Read messages from backend
	go func() {
		for {
			_, msg, err := ws.ReadMessage()
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Error("Error reading message.")
				continue
			}

			parts := strings.SplitN(string(msg), common.MessageSeparator, 2)
			clientKey := parts[0]
			msgString := parts[1]
			if client, ok := m.clients[clientKey]; ok {
				client <- msgString
			} else {
				log.WithFields(log.Fields{
					"key": clientKey,
				}).Warn("Could not find channel for message. Dropping message and sending close to backend.")
				m.sendClose(clientKey)
			}
		}
	}()

	// Write messages to backend
	go func() {
		for {
			message := <-m.messagesToBackend
			err := ws.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Error("Error writing message.")
			}
		}
	}()
}
