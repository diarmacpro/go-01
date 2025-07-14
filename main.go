package main

import (
	"fmt"
	"image/png"
	"net/http"

	"github.com/diarmacpro/go-wa-client/internal/client"
)

var waClient *client.WhatsAppClient

func main() {
	waClient = client.NewClient()

	http.HandleFunc("/qrcode", func(w http.ResponseWriter, r *http.Request) {
		// 1. Connect websocket
		if err := waClient.ConnectSocket(); err != nil {
			http.Error(w, "Gagal koneksi websocket: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// 2. Kirim client hello
		if err := waClient.SendClientHello(); err != nil {
			http.Error(w, "Gagal kirim hello: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// 3. Terima ref dari server
		ref, err := waClient.ReadRef()
		if err != nil {
			http.Error(w, "Gagal terima ref: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// 4. Render QR code dari ref
		img, err := client.RenderQR(ref.Ref)
		if err != nil {
			http.Error(w, "Gagal render QR", http.StatusInternalServerError)
			return
		}

		// 5. Kirim ke browser
		w.Header().Set("Content-Type", "image/png")
		png.Encode(w, img)
	})

	// Static files (index.html dan script.js)
	http.Handle("/", http.FileServer(http.Dir("static")))

	fmt.Println("Server aktif di http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
