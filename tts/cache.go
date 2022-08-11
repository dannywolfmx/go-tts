package tts

import (
	"crypto/sha1"
	"encoding/hex"
)

func hash(text string) string {
	hasher := sha1.New()

	hasher.Write([]byte(text))
	buff := hasher.Sum(nil)

	return hex.EncodeToString(buff)
}
