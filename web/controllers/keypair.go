package controllers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"io"

	"golang.org/x/crypto/argon2"
)

// RandomBytes generates fixed length random bytes
func RandomBytes(bytes int) (data []byte) {
	data = make([]byte, bytes)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		panic(err)
	}
	return
}

// GenerateRSAKey generates a new key pair
func GenerateRSAKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, RSAKeySize)
}

// RSAEncrypt does public key encryption using RSA-OAEP,
// using sha256 as the hash and the label is nil
func RSAEncrypt(pub *rsa.PublicKey, msg []byte, tag []byte) ([]byte, error) {
	return rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		pub, msg, tag,
	)
}

// RSADecrypt does public key decryption
func RSADecrypt(priv *rsa.PrivateKey, msg []byte, tag []byte) ([]byte, error) {
	return rsa.DecryptOAEP(
		sha256.New(),
		rand.Reader,
		priv, msg, tag,
	)
}

// Argon2Key - Automatically choses a decent combination of iterations and memory
func Argon2Key(password []byte, salt []byte, keyLen uint32) []byte {
	return argon2.IDKey(
		password, salt,
		1, 64*1024, 4,
		keyLen,
	)
}

// CFBEncrypter counter block AES encryption
func CFBEncrypter(key []byte, iv []byte) cipher.Stream {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	return cipher.NewCFBEncrypter(block, iv)
}

// CFBDecrypter counter block AES decryption
func CFBDecrypter(key []byte, iv []byte) cipher.Stream {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	return cipher.NewCFBDecrypter(block, iv)
}
