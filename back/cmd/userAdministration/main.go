package main

import (
	"fmt"
	"note/db"
	"os"
	"strings"
	"syscall"

	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
	"golang.org/x/term"
)

func main() {
	db.Init()
	fmt.Println("Select option:")
	fmt.Println("1. Create user")
	fmt.Println("2. Edit user")
	fmt.Println("3. Delete user")
	option := 0
	fmt.Scanf("%d", &option)
	fmap := make(map[int]func())
	fmap[1] = CreateUser
	fmap[2] = EditUser
	fmap[3] = DeleteUser
	fmap[option]()
}

func CreateUser() {
	fmt.Printf("Username: ")
	var username string
	fmt.Scanf("%s", &username)

	fmt.Printf("Password: ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Printf("Error reading password: %v\n", err)
		os.Exit(1)
	}
	password := string(passwordBytes)
	fmt.Println()

	totpKey := genTotp(username)

	err = db.CreateUser(username, password, totpKey)
	if err != nil {
		fmt.Println("Failed to create user:", err)
		os.Exit(1)
	}

	fmt.Println("User created")
}

func EditUser() {
	fmt.Printf("Username to edit: ")
	var username string
	fmt.Scanf("%s", &username)

	fmt.Printf("New password (leave blank to keep unchanged): ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Printf("Error reading password: %v\n", err)
		os.Exit(1)
	}
	newPassword := string(passwordBytes)
	fmt.Println()

	fmt.Printf("Do you want to generate a new TOTP code? (y/N)")
	totpKey := ""
	confirm := ""
	fmt.Scanf("%s", &confirm)
	if strings.ToUpper(confirm) == "Y" {
		totpKey = genTotp(username)
	}

	err = db.EditUser(username, newPassword, totpKey)
	if err != nil {
		fmt.Println("Failed to edit user:", err)
		os.Exit(1)
	}

	fmt.Println("User updated")
}

func DeleteUser() {
	fmt.Printf("Username to delete: ")
	var username string
	fmt.Scanf("%s", &username)

	fmt.Printf("Are you sure you want to delete user '%s'? (y/N): ", username)
	var confirm string
	fmt.Scanf("%s", &confirm)
	if confirm != "y" && confirm != "Y" {
		fmt.Println("Aborted deletion")
		return
	}

	err := db.DeleteUser(username)
	if err != nil {
		fmt.Println("Failed to delete user:", err)
		os.Exit(1)
	}

	fmt.Println("User deleted")
}

func genTotp(username string) string {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "note",
		AccountName: username,
	})
	if err != nil {
		fmt.Println("error while creating totp key:", err)
		os.Exit(1)
	}

	display(key.URL())
	fmt.Println("TOTP URL:", key.URL())

	var passcode string
	fmt.Print("Type the generated code: ")
	fmt.Scanf("%s", &passcode)

	if !totp.Validate(passcode, key.Secret()) {
		fmt.Println("Failed to verify TOTP")
		os.Exit(1)
	}
	return key.Secret()
}

func display(data string) {
	qr, err := qrcode.New(data, qrcode.Low)
	if err != nil {
		fmt.Println("Failed to generate ASCII QR code:", err)
		os.Exit(1)
	}

	fmt.Println(qr.ToString(true))
}
