package utils_encryption_test

import (
	utils_encryption "restful-service/utils/encryption"
	"testing"
)

var secretKey string = "a007c982fc573c32b647a21b12561e05"

func TestEncrypt_ShouldNotReturnSamecipher(t *testing.T) {
	plainText := "hello world"

	cipher1, _ := utils_encryption.Encrypt(plainText, secretKey)
	cipher2, _ := utils_encryption.Encrypt(plainText, secretKey)

	if cipher1 == cipher2 {
		t.Error("Encryption returning same cipher")
	}
}

func TestEncrypt_ShouldNotWorkForInvalidSecretKey(t *testing.T) {
	plainText := "hello world"

	_, err := utils_encryption.Encrypt(plainText, "invalid-secret")

	if err == nil {
		t.Error("Encryption is working with invalid secret")
	}
}

func TestDecrypt_ShouldReturnErrorForEmptyString(t *testing.T) {
	cipher := ""

	_, err := utils_encryption.Decrypt(cipher, secretKey)

	if err.Error() != "unable to decrypt empty string" {
		t.Error("Decryption not returning error for empty stirng")
	}
}

func TestDecrypt_ShouldReturnErrorForInvalidcipher(t *testing.T) {
	cipher := "invalid-cipher"

	_, err := utils_encryption.Decrypt(cipher, secretKey)

	if err == nil {
		t.Error("Decryption not returning error for invalid cipher")
	}
}

func TestDecrypt_ShouldDecryptcipher(t *testing.T) {
	cipher := "VKRsTNEtEOWRNQom3vdd73AIBVrFhbw31eX5KfIeuXk="

	result, err := utils_encryption.Decrypt(cipher, secretKey)

	if err != nil {
		t.Error("Decrypt unable to decrypt cipher")
	}

	if result != "hello world" {
		t.Error("Decrypt unable to decrypt cipher")
	}
}

func FuzzEncryptionDecryption(f *testing.F) {
	f.Add("owh!$%%#")
	f.Fuzz(func(t *testing.T, plainText string) {
		cipherText, err := utils_encryption.Encrypt(plainText, secretKey)
		if err != nil {
			t.Error("Unable to encrypt")
		}

		decryptedText, err := utils_encryption.Decrypt(cipherText, secretKey)
		if err != nil {
			t.Error("Unable to decrypt")
		}

		if plainText != decryptedText {
			t.Error("Encryption and decryption not working")
		}

	})
}
