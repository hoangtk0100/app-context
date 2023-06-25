package token

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	appctx "github.com/hoangtk0100/app-context"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

const (
	minJWTSecretKeySize          = 32
	defaultAccessTokenExpiresIn  = time.Hour * 24 * 7
	defaultRefreshTokenExpiresIn = time.Hour * 24 * 7 * 2
)

var (
	ErrInvalidJWTKeySize     = errors.New(fmt.Sprintf("Invalid key size: must be at least %v characters", minJWTSecretKeySize))
	ErrMissingCustomDuration = errors.New("Duration must be provided for CustomToken")
	ErrTooManyCustomDuration = errors.New("Provide too many durations")
	ErrInvalidTokenType      = errors.New("Invalid token type")
)

// jwtMaker is a JSON Web Token maker
// Symmetric key algorithm to sign the key
type jwtMaker struct {
	id        string
	secretKey string
	*tokenOpt
}

func NewJWTMaker(id string) *jwtMaker {
	return &jwtMaker{
		id:       id,
		tokenOpt: new(tokenOpt),
	}
}

func (maker *jwtMaker) ID() string {
	return maker.id
}

func (maker *jwtMaker) InitFlags() {
	pflag.StringVar(&maker.secretKey,
		"jwt-secret-key",
		"",
		fmt.Sprintf("JWT secret key - Key size >= %d", minJWTSecretKeySize),
	)

	pflag.DurationVar(&maker.accessTokenExpiresIn,
		"access-token-expires-in",
		defaultAccessTokenExpiresIn,
		"Access token expires in duration - Ex: 2h - Default: 168h",
	)

	pflag.DurationVar(&maker.refreshTokenExpiresIn,
		"refresh-token-expires-in",
		defaultRefreshTokenExpiresIn,
		"Refresh token expires in duration - Ex: 2h - Default: 336h",
	)
}

func (maker *jwtMaker) Run(_ appctx.AppContext) error {
	if len(maker.secretKey) < minJWTSecretKeySize {
		return errors.WithStack(ErrInvalidJWTKeySize)
	}

	return nil
}

func (maker *jwtMaker) Stop() error {
	return nil
}

func (opt *tokenOpt) getTokenDuration(tokenType TokenType, duration ...time.Duration) (time.Duration, error) {
	var tokenDuration time.Duration

	switch tokenType {
	case AccessToken:
		tokenDuration = opt.accessTokenExpiresIn
	case RefreshToken:
		tokenDuration = opt.refreshTokenExpiresIn
	case CustomToken:
		if len(duration) == 0 {
			return 0, errors.WithStack(ErrMissingCustomDuration)
		} else if len(duration) > 1 {
			return 0, errors.WithStack(ErrTooManyCustomDuration)
		}

		tokenDuration = duration[0]
	default:
		return 0, errors.WithStack(ErrInvalidTokenType)
	}

	return tokenDuration, nil
}

func (maker *jwtMaker) CreateToken(uid string, tokenType TokenType, duration ...time.Duration) (string, *Payload, error) {
	tokenDuration, err := maker.tokenOpt.getTokenDuration(tokenType, duration...)
	if err != nil {
		return "", nil, err
	}

	payload, err := NewPayload(uid, tokenDuration)
	if err != nil {
		return "", payload, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	return token, payload, errors.WithStack(err)
}

func (maker *jwtMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.WithStack(ErrInvalidToken)
		}

		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, errors.WithStack(ErrExpiredToken)
		}

		return nil, errors.WithStack(ErrInvalidToken)
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, errors.WithStack(ErrInvalidToken)
	}

	return payload, nil
}
