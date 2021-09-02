package main

// (c) 2021 Frank Engelhardt, <frank@f9e.de>

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"os"

	"golang.org/x/crypto/sha3"
	"golang.org/x/term"
)

func main() {
	flag.Parse()
	username := ""
	if len(os.Args) > 1 {
		username = os.Args[1]
	}
	if username == "" {
		fmt.Printf("Please specify a user name: ")
		fmt.Scanln(&username)
	}
	fmt.Printf("Please enter the passsword for user %s: ", username)
	password, _ := term.ReadPassword(0)
	fmt.Printf("\nPlease re-type the passsword: ")
	password2, _ := term.ReadPassword(0)
	if bytes.Compare(password, password2) == 0 {
		data := fmt.Sprintf("%s:%s", username, password)
		b64data := base64.StdEncoding.EncodeToString([]byte(data))
		hash := make([]byte, 64)
		sha3.ShakeSum256(hash, []byte(b64data))
		fmt.Printf("\n%x\n", hash)
	} else {
		fmt.Println("The passwords do not match. Please try again.")
	}
}
