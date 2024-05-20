package randx

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateBase64(length int) (string, error) {
	b := make([]byte, length)
	n, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	str := base64.StdEncoding.EncodeToString(b[:n])
	if len(str) > length {
		str = str[:length]
	}

	return str, nil
}