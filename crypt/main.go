/*
 * Genarate rsa keys.
 */

package crypt

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"os"
)

func GetCerts() (*rsa.PrivateKey, error) {
	reader := rand.Reader
	bitSize := 2048
	return rsa.GenerateKey(reader, bitSize)
}

func SignData(data string, key *rsa.PrivateKey) ([]byte, error) {
	hashed := sha256.Sum256([]byte(data))
	return key.Sign(rand.Reader, hashed[:], crypto.SHA256)
}

func MarshalPublicKey(pubkey rsa.PublicKey) ([]byte, error) {
	asn1Bytes, err := asn1.Marshal(pubkey)
	checkError(err)

	var pemkey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	var buf bytes.Buffer
	err = pem.Encode(&buf, pemkey)
	return buf.Bytes(), err
}
func UnMarshalPublicKey(key []byte) (*rsa.PublicKey, error) {

	pk := &rsa.PublicKey{}
	_, err := asn1.Unmarshal(key, pk)
	checkError(err)

	return pk, err
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
