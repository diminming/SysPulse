package common

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5Calc(s []byte) string {
	ret := md5.Sum(s)
	return hex.EncodeToString(ret[:])
}
