package utils

import (
	"bytes"
	"encoding/base64"
	"image/png"

	qrcode "github.com/skip2/go-qrcode"
)

// GenerateQRCodePNG generates a QR code from data and returns it as a PNG image in base64 format.
// The QR code level is set to Medium (40% error correction).
// Size parameter controls the size of the QR code (256 is recommended for small displays).
func GenerateQRCodePNG(data string, size int) (string, error) {
	// Create a new QR code with the provided data
	qr, err := qrcode.New(data, qrcode.Medium)
	if err != nil {
		return "", err
	}

	// Set the size for the QR code
	// Each module (dot) will be 'size/QR_WIDTH' pixels (approximately)
	qr.DisableBorder = false

	// Create a buffer to write the PNG image
	var buf bytes.Buffer

	// Encode the QR code as PNG to the buffer
	qrImage := qr.Image(256) // 256x256 pixel image
	if err := png.Encode(&buf, qrImage); err != nil {
		return "", err
	}

	// Convert the PNG image to base64 string
	base64String := base64.StdEncoding.EncodeToString(buf.Bytes())
	return base64String, nil
}

// GenerateQRCodePNGWithLevel generates a QR code with a specific error correction level.
// Levels: Low (7%), Medium (15%), High (30%), Highest (50%)
func GenerateQRCodePNGWithLevel(data string, level qrcode.RecoveryLevel, size int) (string, error) {
	// Create a new QR code with the provided data and level
	qr, err := qrcode.New(data, level)
	if err != nil {
		return "", err
	}

	qr.DisableBorder = false

	// Create a buffer to write the PNG image
	var buf bytes.Buffer

	// Encode the QR code as PNG to the buffer
	qrImage := qr.Image(256)
	if err := png.Encode(&buf, qrImage); err != nil {
		return "", err
	}

	// Convert the PNG image to base64 string
	base64String := base64.StdEncoding.EncodeToString(buf.Bytes())
	return base64String, nil
}

// ValidateQRCodeToken checks if the QR code token is valid (you can add custom logic here)
// For now, it simply checks if the token is not empty
func ValidateQRCodeToken(token string) bool {
	return token != ""
}
