package utils_encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

func Encrypt(text string, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	iv := make([]byte, aes.BlockSize)
	_, err = io.ReadFull(rand.Reader, iv)

	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	plainText := pkcs5Padding([]byte(text), block.BlockSize())

	cipherText := make([]byte, len(plainText))
	mode.CryptBlocks(cipherText, plainText)

	combined := append(iv, cipherText...)

	encryptedPassword := base64.StdEncoding.EncodeToString(combined)

	return encryptedPassword, nil
}

func Decrypt(encryptedData string, key string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		fmt.Println("Error decoding:", err)
		return "", err
	}

	// Create AES cipher block
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		fmt.Println("Error creating cipher block:", err)
		return "", err
	}

	// Determine the initialization vector (IV) size based on the block size
	iv := ciphertext[:aes.BlockSize]

	// Use CBC mode for decryption
	mode := cipher.NewCBCDecrypter(block, iv)

	// Decrypt the data
	plainText := make([]byte, len(ciphertext)-aes.BlockSize)
	mode.CryptBlocks(plainText, ciphertext[aes.BlockSize:])

	// Remove padding
	plainText = pkcs5Unpad(plainText)

	return string(plainText), nil
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs5Unpad(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}
