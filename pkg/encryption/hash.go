package encryption

import (
	"crypto/md5"
	"encoding/hex"
)

func Hash(data []byte) string {
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}
