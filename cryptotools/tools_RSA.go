package cryptotools

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"os"
)

func GenerateRSAKey(bits int) {

	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		panic(err)
	}

	X509PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)
	privateFile, err := os.Create("private.pem")
	if err != nil {
		panic(err)
	}
	defer privateFile.Close()

	privateBlock := pem.Block{Type: "RSA Private Key", Bytes: X509PrivateKey}
	pem.Encode(privateFile, &privateBlock)
	publicKey := privateKey.PublicKey

	X509PublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		panic(err)
	}

	publicFile, err := os.Create("public.pem")
	if err != nil {
		panic(err)
	}
	defer publicFile.Close()

	publicBlock := pem.Block{Type: "RSA Public Key", Bytes: X509PublicKey}
	pem.Encode(publicFile, &publicBlock)

}

func RSA_Encrypt(plainText []byte, path string) []byte {

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)

	block, _ := pem.Decode(buf)

	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	publicKey := publicKeyInterface.(*rsa.PublicKey)
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)
	if err != nil {
		panic(err)
	}

	return cipherText

}

func RSA_Decrypt(cipherText []byte, path string) []byte {

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)

	block, _ := pem.Decode(buf)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	plainText, _ := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherText)

	return plainText

}

func RSA_DigitalSignature(msg []byte, path string) []byte {

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)

	block, _ := pem.Decode(buf)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	msghash := sha512.Sum512(msg)
	signbyt, _ := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA512, msghash[:])

	return signbyt

}

func RSA_DigitalSignatureVerify(msg []byte, sign []byte, path string) bool {

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)

	block, _ := pem.Decode(buf)

	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	publicKey := publicKeyInterface.(*rsa.PublicKey)

	msghash := sha512.Sum512(msg)
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA512, msghash[:], sign)

	if err != nil {
		return false
	}
	return true

}

func GetKey(path string) string {

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)

	return string(buf)

}

func EncryptRSAToString(msg string) string {

	cmsg := RSA_Encrypt([]byte(msg), "public.pem")
	cmsgb64 := base64.StdEncoding.EncodeToString(cmsg)
	return cmsgb64

}

func DecryptRSAToString(cmsgb64 string) string {

	cmsg, _ := base64.StdEncoding.DecodeString(cmsgb64)
	return string(RSA_Decrypt(cmsg, "private.pem"))

}

func DigitalSignature(msg string) string {

	msgin := RSA_DigitalSignature([]byte(msg), "private.pem")
	msginb64 := base64.StdEncoding.EncodeToString(msgin)
	return msginb64

}

func DigitalSignatureVerify(msg string, sign string) bool {

	msgin, _ := base64.StdEncoding.DecodeString(sign)
	return RSA_DigitalSignatureVerify([]byte(msg), msgin, "public.pem")

}

func init() {

	GenerateRSAKey(2048)

}
