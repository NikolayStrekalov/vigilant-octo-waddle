package sign

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func Sign(b []byte, key string) (string, error) {
	h := sha256.New()
	_, err := h.Write(b)
	if err != nil {
		return "", fmt.Errorf("sign error1: %w", err)
	}
	_, err = h.Write([]byte(key))
	if err != nil {
		return "", fmt.Errorf("sign error2: %w", err)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
