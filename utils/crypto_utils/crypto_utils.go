package crypto_utils


import (
	"encoding/hex"
	"crypto/md5"
)

func GetMd5(input string)string  {
	hash := md5.New()
	defer hash.Reset()
	hash.Write([]byte(input))
	return hex.EncodeToString(hash.Sum(nil))
}