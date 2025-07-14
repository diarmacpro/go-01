package client

import (
	"log"

	"github.com/gorilla/websocket"
)

type WhatsAppClient struct {
	conn *websocket.Conn
}

func NewClient() *WhatsAppClient {
	return &WhatsAppClient{}
}

func (w *WhatsAppClient) GetDummyQR() string {
	return "WADUMMY-QR-STRING-123456"
}

func (w *WhatsAppClient) ConnectSocket() error {
	url := "wss://web.whatsapp.com/ws"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	log.Println("WebSocket connected")
	w.conn = conn
	return nil
}
