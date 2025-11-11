package pkg

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashToSha256(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:8])
}
