package service

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// 生成和验证token
type Claims struct {
	UserID int64
	jwt.StandardClaims
}

var jwtKey = []byte("tokenkey")

// 生成token
func GenerateToken(id int64) (string, error) {
	expirationTime := time.Now().UTC().Add(24 * time.Hour)
	claims := &Claims{
		UserID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// 验证token
func VerifyToken(tknStr string) (*Claims, error) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, fmt.Errorf("invalid JWT signature")
		}
		return nil, fmt.Errorf("could not parse JWT token")
	}

	if !tkn.Valid {
		return nil, fmt.Errorf("invalid JWT token")
	}

	return claims, nil
}
