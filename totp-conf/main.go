package main

import (
	"fmt"
	"os"

	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s USER_NAME\n", os.Args[0])
		os.Exit(1)
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "note",
		AccountName: os.Args[1],
	})

	if err != nil {
		fmt.Print("error while creating totp key")
		os.Exit(1)
	}

	display(key.URL())

	var passcode string
	fmt.Print("Type the generated code: ")
	fmt.Scanf("%s", &passcode)

	if !totp.Validate(passcode, key.Secret()) {
		fmt.Println("Failed to verify TOTP")
		os.Exit(1)
	}
	fmt.Println("Add to config.toml:")
	fmt.Printf("totp=%s\n", key.Secret())
}

func display(data string) {
	qr, err := qrcode.New(data, qrcode.Low)
	if err != nil {
		fmt.Println("Failed to generate ASCII QR code:", err)
		os.Exit(1)
	}

	fmt.Println(qr.ToString(true))
}
