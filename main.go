package main

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/validate", validateHandler)
	fmt.Println("Server is running on :8080")
	http.ListenAndServe(":8080", nil)
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form data
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit for file size
	if err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	// Get the uploaded file
	file, _, err := r.FormFile("ticketData")
	if err != nil {
		http.Error(w, "Failed to retrieve the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read the content of the uploaded file
	ticketData, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read the uploaded file", http.StatusBadRequest)
		return
	}

	// Parse the public key from the command-line argument (Replace with your key file path)
	keyfileStr := "./public_key.pem"
	keyfile, err := ioutil.ReadFile(keyfileStr)
	if err != nil {
		http.Error(w, "Failed to read public key file", http.StatusInternalServerError)
		return
	}

	pemBlock, _ := pem.Decode(keyfile)
	key, err := x509.ParsePKIXPublicKey(pemBlock.Bytes)
	if err != nil {
		http.Error(w, "Failed to parse public key", http.StatusInternalServerError)
		return
	}
	pubKey := key.(*rsa.PublicKey)

	ticketArr := strings.Split(string(ticketData), ".")
	fmt.Printf("Payload: %s\n", ticketArr[0])
	fmt.Printf("Signature: %s\n", ticketArr[1])

	payload, err := base64.StdEncoding.DecodeString(ticketArr[0])
	fmt.Println(err)
	if err != nil {
		http.Error(w, "Failed to decode payload", http.StatusBadRequest)
		return
	}
	sig, err := base64.StdEncoding.DecodeString(ticketArr[1])
	if err != nil {
		http.Error(w, "Failed to decode signature", http.StatusBadRequest)
		return
	}

	hash := sha256.Sum256(payload)

	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hash[:], sig)
	if err != nil {
		http.Error(w, "Signature verification failed", http.StatusUnauthorized)
		return
	}

	// Validation successful, respond with success message
	fmt.Fprintln(w, "Validation successful")
}
