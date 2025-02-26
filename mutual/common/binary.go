package common

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"hash/fnv"

	"go.uber.org/zap"
)

func MD5Calc(s []byte) string {
	ret := md5.Sum(s)
	return hex.EncodeToString(ret[:])
}

func GetHash(obj any) string {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(obj)
	zap.L().Error("GetHash", zap.Error(err))

	h := fnv.New32a()
	h.Write(buf.Bytes())
	return hex.EncodeToString(h.Sum(nil))
}

func GetStringHash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
