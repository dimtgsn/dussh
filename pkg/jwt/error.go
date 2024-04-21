package jwt

import "errors"

var (
	ErrEmptySecretKey          = errors.New("empty secret key")
	ErrCookieNotFound          = errors.New("cookie not found")
	ErrInvalidToken            = errors.New("invalid token")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method: %v")
	ErrTokenClaimsNotFound     = errors.New("token claims not found")
)
