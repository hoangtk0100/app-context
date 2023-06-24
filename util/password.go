package util

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"golang.org/x/crypto/bcrypt"
)

const (
	defaultSeparator = "."
)

func getFormattedPassword(format string, args ...interface{}) (string, error) {
	format = strings.TrimSpace(format)
	if format != "" {
		val := fmt.Sprintf(format, args...)
		if strings.Contains(val, "%!") {
			return "", errors.New("Invalid password format")
		}

		return val, nil
	}

	var password strings.Builder

	for _, val := range args {
		password.WriteString(fmt.Sprintf("%s%v", defaultSeparator, val))
	}

	return password.String()[1:], nil
}

func HashPassword(format string, args ...interface{}) (string, error) {
	password, err := getFormattedPassword(format, args...)
	if err != nil {
		return "", errors.WithStack(err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return string(hashedPassword), nil
}

func CheckPassword(hashedPassword string, format string, args ...interface{}) error {
	password, err := getFormattedPassword(format, args...)
	if err != nil {
		return errors.WithStack(err)
	}

	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
