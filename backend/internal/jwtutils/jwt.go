package jwtutils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claim struct {
	IDUser int32    `json:"id_user"`
	Roles  []string `json:"roles"`
	Email  string   `json:"email"`
	NIK    string   `json:"nik"`
	jwt.RegisteredClaims
}

type JWT struct {
	secret string
}

func New(secret string) JWT {
	return JWT{secret: secret}
}

func (j *JWT) Encode(claim Claim) (string, error) {
	return j.EncodeWithTTL(claim, 30*time.Minute)
}

func (j *JWT) EncodeWithTTL(claim Claim, ttl time.Duration) (string, error) {
	claim.ExpiresAt = jwt.NewNumericDate(time.Now().Add(ttl))
	claim.IssuedAt = jwt.NewNumericDate(time.Now())

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claim)
	return token.SignedString([]byte(j.secret))
}

func (j *JWT) Decode(token string) (*Claim, error) {
	claims := &Claim{}
	parsed, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !parsed.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
