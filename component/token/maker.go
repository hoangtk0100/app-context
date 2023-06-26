package token

import (
	"time"
)

type TokenType int

const (
	AccessToken TokenType = iota
	RefreshToken
	CustomToken
)

type TokenMaker interface {
	// CreateToken creates a new token for a specific uid and duration
	// For TokenType CustomToken, duration must be provided explicitly.
	CreateToken(tokenType TokenType, uid string, duration ...time.Duration) (string, *Payload, error)

	// VerifyToken checks if a token is valid or not
	VerifyToken(token string) (*Payload, error)
}
