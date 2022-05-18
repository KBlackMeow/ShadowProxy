package cryptotools

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"time"
)

var Salt = fmt.Sprint(time.Now().UnixMilli() % 556639)

func Hash_SHA(str string) string {

	str = str + Salt
	hashbyt := sha256.Sum256([]byte(str))
	return hex.EncodeToString(hashbyt[:])

}

func Hash_SHA512(str string) string {

	str = str + Salt
	str = Hash_SHA(str)
	hashbyt := sha512.Sum512([]byte(str))
	return hex.EncodeToString(hashbyt[:])

}

func Hash_SHA256(str string) string {

	str = str + Salt
	str = Hash_SHA(str)
	hashbyt := sha256.Sum256([]byte(str))
	return hex.EncodeToString(hashbyt[:])

}

func Hash_MD5(str string) string {

	str = str + Salt
	str = Hash_SHA(str)
	hashbyt := md5.Sum([]byte(str))
	return hex.EncodeToString(hashbyt[:])

}
