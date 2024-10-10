package panelist

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type MessageType string

const (
	MessageTypeInit MessageType = "init"
)

type panelist struct {
	config        Config
	websocketConn *websocket.Conn
}

type Panelist interface {
}

type ServerConfig struct {
	Host string
	Port int
}

type Config struct {
	ClientName string
	Server     ServerConfig
	Pages      []string
}

func New(config Config) Panelist {
	return &panelist{}
}

func (p *panelist) Init() error {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf(`ws://%s:%d/ws`, p.config.Server.Host, p.config.Server.Port), nil)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}
	p.websocketConn = conn
	p.sendMessage(MessageTypeInit, map[string]interface{}{"client_name": p.config.ClientName})
	return nil
}

func (p *panelist) sendMessage(messageType MessageType, message interface{}) error {
	return p.websocketConn.WriteJSON(map[string]interface{}{
		"type": messageType,
		"data": message,
	})
}
