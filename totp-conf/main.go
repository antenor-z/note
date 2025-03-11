package main

import (
	"bytes"
	"fmt"
	"image/png"
	"os"

	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s DOMAIN_NAME USER_NAME\n", os.Args[0])
		os.Exit(1)
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      os.Args[1],
		AccountName: os.Args[2],
	})

	if err != nil {
		fmt.Print("error while creating totp key")
		return
	}

	var buf bytes.Buffer
	img, err := key.Image(50, 50)
	if err != nil {
		fmt.Print("error while creating totp qr code image")
		return
	}
	png.Encode(&buf, img)

	display(buf.Bytes())

	var passcode string
	fmt.Print("Type the generated code: ")
	fmt.Scanf("%s", passcode)

	if !totp.Validate(passcode, key.Secret()) {
		fmt.Println("Failed to verify TOTP")
		os.Exit(1)
	}
	fmt.Println("Add to config.toml:")
	fmt.Printf("totpKey=%s", key.Secret())
}

func display(data []byte) {
	qr, err := qrcode.New(string(data), qrcode.Low)
	if err != nil {
		fmt.Println("Failed to generate ASCII QR code:", err)
		return
	}

	fmt.Println(qr.ToString(true))
}
