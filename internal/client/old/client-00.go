package client

type WhatsAppClient struct{}

func NewClient() *WhatsAppClient {
	return &WhatsAppClient{}
}

// Sementara: QR dummy (ganti nanti dengan QR asli dari protokol WhatsApp)
func (w *WhatsAppClient) GetDummyQR() string {
	return "WADUMMY-QR-STRING-123456"
}
