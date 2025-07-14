package client

import (
	"encoding/json"
	"log"

	"github.com/flynn/noise"
	"github.com/gorilla/websocket"
)

type WhatsAppClient struct {
	conn *websocket.Conn
	hs   *noise.HandshakeState
}

type RefResponse struct {
	Ref       string `json:"ref"`
	PublicKey string `json:"publicKey"`
	ClientID  string `json:"clientID"`
	TTL       int    `json:"ttl"`
}

func NewClient() *WhatsAppClient {
	return &WhatsAppClient{}
}

// Membuka koneksi ke server WhatsApp Web
func (w *WhatsAppClient) ConnectSocket() error {
	url := "wss://web.whatsapp.com/ws/chat"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	log.Println("WebSocket connected")
	w.conn = conn
	return nil
}

// Mengirim Client Hello (Noise handshake tahap 1)
func (w *WhatsAppClient) SendClientHello() error {
	config := noise.Config{
		Pattern:   noise.HandshakeXX,
		Initiator: true,
		Prologue:  []byte("Noise_XX_25519_AESGCM_SHA256"),
	}
	hs, err := noise.NewHandshakeState(config)
	if err != nil {
		return err
	}

	w.hs = hs // simpan state untuk lanjutkan nanti

	msg1, _, _, err := hs.WriteMessage(nil, nil)
	if err != nil {
		return err
	}

	finalPayload := append([]byte{0x00, 0x01}, msg1...)

	err = w.conn.WriteMessage(websocket.BinaryMessage, finalPayload)
	if err != nil {
		return err
	}

	log.Println("Client Hello sent")
	return nil
}

// Membaca Server Hello + pesan 'ref' untuk login (QR)
func (w *WhatsAppClient) ReadRef() (*RefResponse, error) {
	_, msg, err := w.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	// Lanjutkan Noise handshake (read message 2)
	plaintext, _, _, err := w.hs.ReadMessage(nil, msg)
	if err != nil {
		log.Println("Noise read failed:", err)
		return nil, err
	}

	log.Println("Plaintext from WA:", string(plaintext))

	var ref RefResponse
	if err := json.Unmarshal(plaintext, &ref); err != nil {
		log.Println("JSON parse failed:", err)
		return nil, err
	}

	return &ref, nil
}
