package token

import (
	"fmt"
	"time"

	appctx "github.com/hoangtk0100/app-context"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

const (
	pasetoSymmetricKeySize = chacha20poly1305.KeySize
)

var (
	ErrInvalidPasetoKeySize = errors.New(fmt.Sprintf("Invalid key size: must be exactly %d charaters", pasetoSymmetricKeySize))
)

type tokenOpt struct {
	accessTokenExpiresIn  time.Duration
	refreshTokenExpiresIn time.Duration
}

// PasetoMaker is a PASETO token maker
// use symmetric encryption to encrypt the token payload
type pasetoMaker struct {
	id           string
	paseto       *paseto.V2
	symmetricKey []byte
	*tokenOpt
}

func NewPasetoMaker(id string) *pasetoMaker {
	return &pasetoMaker{
		id:       id,
		paseto:   paseto.NewV2(),
		tokenOpt: new(tokenOpt),
	}
}

func (maker *pasetoMaker) ID() string {
	return maker.id
}

func (maker *pasetoMaker) InitFlags() {
	pflag.BytesHexVar(&maker.symmetricKey,
		"paseto-symmetric-key",
		[]byte(""),
		fmt.Sprintf("PASETO symmetric key - Key size: %d", pasetoSymmetricKeySize),
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

func (maker *pasetoMaker) Run(_ appctx.AppContext) error {
	if len(maker.symmetricKey) != pasetoSymmetricKeySize {
		return errors.WithStack(ErrInvalidPasetoKeySize)
	}

	return nil
}

func (maker *pasetoMaker) Stop() error {
	return nil
}

func (maker *pasetoMaker) CreateToken(tokenType TokenType, uid string, duration ...time.Duration) (string, *Payload, error) {
	tokenDuration, err := maker.tokenOpt.getTokenDuration(tokenType, duration...)
	if err != nil {
		return "", nil, err
	}

	payload, err := NewPayload(uid, tokenDuration)
	if err != nil {
		return "", payload, err
	}

	token, err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	return token, payload, errors.WithStack(err)
}

func (maker *pasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, errors.WithStack(ErrInvalidToken)
	}

	err = payload.Valid()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return payload, nil
}
