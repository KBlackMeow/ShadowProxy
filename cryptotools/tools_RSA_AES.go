package cryptotools

import (
	"bytes"
	"crypto/aes"
)

func RSA_AES_encode(pubkey string, key string, msg []byte) []byte {

	if pubkey == "" {
		ckey := RSA_Encrypt([]byte(key), "public.pem")
		cmsg := Ase256Encode(msg, key, "1234567890123456", aes.BlockSize)
		var buff bytes.Buffer
		buff.Write(ckey)
		buff.Write(cmsg)
		return buff.Bytes()
	} else {
		ckey := RSA_Encode([]byte(key), pubkey)
		cmsg := Ase256Encode(msg, key, "1234567890123456", aes.BlockSize)
		var buff bytes.Buffer
		buff.Write(ckey)
		buff.Write(cmsg)
		return buff.Bytes()
	}
}

func RSA_AES_decode(prikey string, cmsg []byte) ([]byte, string) {
	if prikey == "" {
		key := RSA_Decrypt(cmsg[:256], "private.pem")
		msg := Ase256Decode(cmsg[256:], string(key), "1234567890123456")
		return msg, string(key)
	} else {
		key := RSA_Decode(cmsg[:256], prikey)
		msg := Ase256Decode(cmsg[256:], string(key), "1234567890123456")
		return msg, string(key)
	}
}
