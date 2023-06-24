package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password, err := RandomString(6)
	require.NoError(t, err)
	require.NotEmpty(t, password)

	hashedPassword1, err := HashPassword("", password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)

	err = CheckPassword(hashedPassword1, "", password)
	fmt.Println(password)
	fmt.Println(hashedPassword1)
	require.NoError(t, err)

	wrongPassword, err := RandomString(6)
	require.NoError(t, err)
	require.NotEmpty(t, wrongPassword)

	err = CheckPassword(hashedPassword1, "", wrongPassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashedPassword2, err := HashPassword("", password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)
	require.NotEqual(t, hashedPassword1, hashedPassword2)

	_, err = HashPassword("%v.%v", password)
	require.Error(t, err)

	salt, err := RandomString(6)
	require.NoError(t, err)
	require.NotEmpty(t, salt)

	hashedPassword3, err := HashPassword("%v.%v", salt, password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword3)
	fmt.Println(hashedPassword3)

	err = CheckPassword(hashedPassword3, "", salt, password)
	require.NoError(t, err)
}
