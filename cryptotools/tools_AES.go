package cryptotools

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
)

func GenAesKey() string {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		panic(err)
	}
	base64Str := base64.StdEncoding.EncodeToString(key)

	return base64Str[:32]
}
func Ase256Encode(plainbyte []byte, key string, iv string, blockSize int) []byte {
	bKey := []byte(key)
	bIV := []byte(iv)
	bPlain := PKCS5Padding(plainbyte, blockSize, len(plainbyte))
	block, err := aes.NewCipher(bKey)
	if err != nil {
		panic(err)
	}
	cipherbyte := make([]byte, len(bPlain))
	mode := cipher.NewCBCEncrypter(block, bIV)
	mode.CryptBlocks(cipherbyte, bPlain)
	return cipherbyte
}

func Ase256Decode(cipherByte []byte, encKey string, iv string) []byte {
	bKey := []byte(encKey)
	bIV := []byte(iv)
	block, err := aes.NewCipher(bKey)
	if err != nil {
		panic(err)
	}

	mode := cipher.NewCBCDecrypter(block, bIV)
	mode.CryptBlocks(cipherByte, cipherByte)
	cipherByte = PKCS5UnPadding(cipherByte)
	return cipherByte
}

func Ase256EncodeHex(plaintext string, key string, iv string, blockSize int) string {
	bKey := []byte(key)
	bIV := []byte(iv)
	bPlaintext := PKCS5Padding([]byte(plaintext), blockSize, len(plaintext))
	block, err := aes.NewCipher(bKey)
	if err != nil {
		panic(err)
	}
	ciphertext := make([]byte, len(bPlaintext))
	mode := cipher.NewCBCEncrypter(block, bIV)
	mode.CryptBlocks(ciphertext, bPlaintext)
	return hex.EncodeToString(ciphertext)
}

func Ase256DecodeHex(cipherText string, encKey string, iv string) (decryptedString string) {
	bKey := []byte(encKey)
	bIV := []byte(iv)
	cipherTextDecoded, err := hex.DecodeString(cipherText)
	if err != nil {
		panic(err)
	}

	block, err := aes.NewCipher(bKey)
	if err != nil {
		panic(err)
	}

	mode := cipher.NewCBCDecrypter(block, bIV)
	mode.CryptBlocks([]byte(cipherTextDecoded), []byte(cipherTextDecoded))
	return string(cipherTextDecoded)
}

func PKCS5Padding(ciphertext []byte, blockSize int, after int) []byte {
	padding := (blockSize - len(ciphertext)%blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
