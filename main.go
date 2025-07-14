package main

import (
	"fmt"
	"image/png"
	"log"
	"net/http"

	"github.com/diarmacpro/go-wa-client/internal/client"
)

var waClient *client.WhatsAppClient

func main() {
	waClient = client.NewClient()

	http.HandleFunc("/qrcode", func(w http.ResponseWriter, r *http.Request) {
		// Langkah 1: Koneksi
		err := waClient.ConnectSocket()
		if err != nil {
			http.Error(w, "Gagal koneksi WebSocket", http.StatusInternalServerError)
			log.Println("WS Error:", err)
			return
		}

		// Langkah 2: Kirim handshake awal
		err = waClient.SendClientHello()
		if err != nil {
			http.Error(w, "Gagal kirim client hello", http.StatusInternalServerError)
			log.Println("Handshake Error:", err)
			return
		}

		// Langkah 3: Terima data awal dari server
		refData, err := waClient.ReadRef()
		if err != nil {
			http.Error(w, "Gagal ambil ref WA", http.StatusInternalServerError)
			log.Println("Read Ref Error:", err)
			return
		}

		log.Println("Dapat REF:", refData.Ref)

		// Langkah 4: Buat QR string dari ref
		qrStr := fmt.Sprintf("wa://login?ref=%s&id=%s&token=%s",
			refData.Ref, refData.ClientID, refData.PublicKey)

		img, err := client.RenderQR(qrStr)
		if err != nil {
			http.Error(w, "Gagal render QR", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "image/png")
		png.Encode(w, img)
	})

	// Static file untuk index.html
	http.Handle("/", http.FileServer(http.Dir("static")))

	fmt.Println("Server aktif di http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
