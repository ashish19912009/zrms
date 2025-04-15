package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func main() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	// Ensure certs folder exists
	if err := os.MkdirAll("certs", os.ModePerm); err != nil {
		panic(err)
	}

	// Save private key
	privateFile, err := os.Create("../certs/private_key.pem")
	if err != nil {
		panic(err)
	}
	defer privateFile.Close()

	privatePEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)
	privateFile.Write(privatePEM)

	// Save public key
	publicFile, err := os.Create("../certs/public_key.pem")
	if err != nil {
		panic(err)
	}
	defer publicFile.Close()

	pubASN1, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		panic(err)
	}

	publicPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: pubASN1,
		},
	)
	publicFile.Write(publicPEM)

	fmt.Println("RSA key pair generated successfully in ./certs/")
}
