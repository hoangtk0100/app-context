package util

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/pkg/errors"
)

func RandomString(length int) (string, error) {
	var sb = make([]byte, length)

	_, err := rand.Read(sb)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return hex.EncodeToString(sb), nil
}
