package client

import (
	"image"
	"bytes"
	"github.com/skip2/go-qrcode"
	"image/png"
)

func RenderQR(text string) (image.Image, error) {
	// Encode ke PNG sebagai []byte
	pngData, err := qrcode.Encode(text, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}

	// Decode kembali ke image.Image
	img, err := png.Decode(bytes.NewReader(pngData))
	if err != nil {
		return nil, err
	}

	return img, nil
}
