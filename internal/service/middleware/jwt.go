package middleware

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type JwtTokenService interface {
	Create(session_id string, tokenExpTime int64) (string, error)
	Validate(tokenString string, expectedSessionId string) (*JwtCsrfClaims, error)
	ParseSecretGetter(token *jwt.Token) (interface{}, error)
}

type JwtToken struct {
	Secret []byte
}

func NewJwtToken(secret string) (JwtTokenService, error) {
	return &JwtToken{
		Secret: []byte(secret),
	}, nil
}

type JwtCsrfClaims struct {
	SessionID string `json:"sid"`
	jwt.StandardClaims
}

func (tk *JwtToken) Create(session_id string, tokenExpTime int64) (string, error) {
	data := JwtCsrfClaims{
		SessionID: session_id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpTime,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	return token.SignedString(tk.Secret)
}

func (tk *JwtToken) Validate(tokenString string, expectedSessionId string) (*JwtCsrfClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtCsrfClaims{}, tk.ParseSecretGetter)
	if err != nil {
		return nil, errors.New("token parse error")
	}

	claims, ok := token.Claims.(*JwtCsrfClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token invalid")
	}

	// Проверка срока действия (дополнительно)
	if claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("token expired")
	}

	if claims.SessionID != expectedSessionId {
		return nil, errors.New("token invalid")
	}

	return claims, nil
}

func (tk *JwtToken) ParseSecretGetter(token *jwt.Token) (interface{}, error) {
	method, ok := token.Method.(*jwt.SigningMethodHMAC)
	if !ok || method.Alg() != "HS256" {
		return nil, errors.New("bad sign method")
	}
	return tk.Secret, nil
}
