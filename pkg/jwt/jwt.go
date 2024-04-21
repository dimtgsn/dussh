package jwt

import (
	"dussh/internal/utils/bytesconv"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"math/rand"
	"strings"
	"time"
)

// TokenManager represents jwt manager.
type TokenManager interface {
	NewAccessToken(userID int64, email string, roleID int) (*Token, error)
	NewRefreshToken() (*Token, error)
}

// UserClaims represents jwt token payloads fields.
type UserClaims struct {
	ID    int64
	Email string
	Role  int
	Exp   int64
}

// TokenPair represents a pair of access and refresh jwt tokens.
type TokenPair struct {
	AccessToken  *Token
	RefreshToken *Token
}

type Token struct {
	Token string
	TTL   int // time in second
}

type tokenManager struct {
	secretKey       string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

// NewManager creates a new token manager.
func NewManager(
	secretKey string,
	accessTokenTTL,
	refreshTokenTTL time.Duration,
) (TokenManager, error) {
	if secretKey == "" {
		return nil, ErrEmptySecretKey
	}

	return &tokenManager{
		secretKey:       secretKey,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}, nil
}

// NewAccessToken creates new access JWT token for given user.
func (tm *tokenManager) NewAccessToken(userID int64, email string, roleID int) (*Token, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = userID
	claims["email"] = email
	claims["role"] = roleID
	claims["exp"] = time.Now().Add(tm.accessTokenTTL).Unix()

	tokenString, err := token.SignedString(bytesconv.StringToBytes(tm.secretKey))
	if err != nil {
		return nil, err
	}

	return &Token{
		Token: tokenString,
		TTL:   int(tm.accessTokenTTL.Seconds()),
	}, nil
}

// NewRefreshToken creates new refresh JWT token.
func (tm *tokenManager) NewRefreshToken() (*Token, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return nil, err
	}

	return &Token{
		Token: fmt.Sprintf("%x", b),
		TTL:   int(tm.refreshTokenTTL.Seconds()),
	}, nil
}

// GetToken returns JWT token.
func GetToken(secretKey, tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf(ErrUnexpectedSigningMethod.Error(), token.Header["alg"])
		}

		return bytesconv.StringToBytes(secretKey), nil
	})

	if err != nil || token == nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	_, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrTokenClaimsNotFound
	}

	return token, nil
}

// GetRefreshToken returns refresh token from cookie.
func GetRefreshToken(c *gin.Context) (string, error) {
	token, err := c.Cookie("RefreshToken")
	if err != nil || token == "" {
		return "", ErrCookieNotFound
	}

	return token, nil
}

// RetrieveJwtToken retrieves and parses a JWT token from a cookie
func RetrieveJwtToken(token *jwt.Token) (*UserClaims, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrTokenClaimsNotFound
	}

	userClaims := &UserClaims{
		ID:    int64(claims["id"].(float64)),
		Email: claims["email"].(string),
		Role:  int(claims["role"].(float64)),
		Exp:   int64(claims["exp"].(float64)),
	}

	return userClaims, nil
}

func ExtractBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("bad header value given")
	}

	jwtToken := strings.Split(header, " ")
	if len(jwtToken) != 2 {
		return "", errors.New("incorrectly formatted authorization header")
	}

	return jwtToken[1], nil
}
