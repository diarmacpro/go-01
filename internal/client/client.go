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

// Langkah 1: Koneksi ke WebSocket WhatsApp
func (w *WhatsAppClient) ConnectSocket() error {
	url := "wss://web.whatsapp.com/ws/chat"
	log.Println("[ConnectSocket] Connecting to:", url)

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Println("[ConnectSocket] Error:", err)
		return err
	}

	log.Println("[ConnectSocket] WebSocket connected")
	w.conn = conn
	return nil
}

// Langkah 2: Kirim Client Hello (Noise Handshake Step 1)
func (w *WhatsAppClient) SendClientHello() error {
	log.Println("[SendClientHello] Preparing Noise handshake")

	config := noise.Config{
		Pattern:   noise.HandshakeXX,
		Initiator: true,
		Prologue:  []byte("Noise_XX_25519_AESGCM_SHA256"),
	}

	hs, err := noise.NewHandshakeState(config)
	if err != nil {
		log.Println("[SendClientHello] Failed to create handshake state:", err)
		return err
	}

	w.hs = hs // simpan handshake state untuk tahap selanjutnya

	msg1, _, _, err := hs.WriteMessage(nil, nil)
	if err != nil {
		log.Println("[SendClientHello] Failed to write message:", err)
		return err
	}

	finalPayload := append([]byte{0x00, 0x01}, msg1...)

	log.Printf("[SendClientHello] Sending client hello (%d bytes)", len(finalPayload))
	err = w.conn.WriteMessage(websocket.BinaryMessage, finalPayload)
	if err != nil {
		log.Println("[SendClientHello] Error sending client hello:", err)
		return err
	}

	log.Println("[SendClientHello] Client Hello sent successfully")
	return nil
}

// Langkah 3: Terima dan parsing pesan ref dari server
func (w *WhatsAppClient) ReadRef() (*RefResponse, error) {
	log.Println("[ReadRef] Waiting for server response...")

	_, msg, err := w.conn.ReadMessage()
	if err != nil {
		log.Println("[ReadRef] Failed to read message:", err)
		return nil, err
	}

	log.Printf("[ReadRef] Received raw message (%d bytes)", len(msg))
	log.Println("[ReadRef] Raw hex:", msg)

	plaintext, _, _, err := w.hs.ReadMessage(nil, msg)
	if err != nil {
		log.Println("[ReadRef] Noise ReadMessage failed:", err)
		return nil, err
	}

	log.Println("[ReadRef] Decrypted plaintext:")
	log.Println(string(plaintext))

	var ref RefResponse
	if err := json.Unmarshal(plaintext, &ref); err != nil {
		log.Println("[ReadRef] JSON parse failed:", err)
		return nil, err
	}

	log.Println("[ReadRef] Parsed ref:", ref.Ref)
	log.Println("[ReadRef] ClientID:", ref.ClientID)
	log.Println("[ReadRef] PublicKey:", ref.PublicKey)

	return &ref, nil
}
