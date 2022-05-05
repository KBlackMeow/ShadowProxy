package cryptotools

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
)

var Salt = fmt.Sprint(time.Now().UnixMilli() % 556639)

func Md5_32(str string) string {
	// logger.Log("Salt", Salt)
	str = str + Salt
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
