package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"

	"code.google.com/p/go.crypto/bcrypt"
)

func main() {
	log.SetFlags(0)
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s username\nPassword is read from stdin.\n"+
			"Password would be stripped from spaces.", os.Args[0])
	}
	username := os.Args[1]
	if strings.Contains(username, ":") {
		log.Fatal("username cannot have colon in it")
	}
	rd := bufio.NewReader(os.Stdin)

	password, err := rd.ReadBytes('\n')
	if err != nil {
		log.Fatal(err)
	}
	password = bytes.TrimSpace(password)
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s:%s\n", username, base64.StdEncoding.EncodeToString(hash))
}
