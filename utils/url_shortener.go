package utils

import (
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/itchyny/base58-go"
)

// Encodes an ID in SHA256
func sha256encoded(input string) ([]byte, error) {
	algorithm := sha256.New()
	_, err := algorithm.Write([]byte(input))
	if err != nil {
		return []byte{}, err
	}
	return algorithm.Sum(nil), nil
}

// Encodes the ID SHA256 encoded in base58
func base58Encoded(bytes []byte) (string, error) {
	encoding := base58.BitcoinEncoding
	encoded, err := encoding.Encode(bytes)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}

// GenerateShortLink generates a short code form a link
func GenerateShortLink(link string, key string) (string, error) {
	urlHashBytes, err := sha256encoded(link + key)
	if err != nil {
		return "", err
	}

	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
	finalString, err := base58Encoded([]byte(fmt.Sprintf("%d", generatedNumber)))
	if err != nil {
		return "", err
	}
	return finalString[:8], nil
}
