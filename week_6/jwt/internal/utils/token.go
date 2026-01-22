package utils

import (
	"time"

	"github.com/airats/micro_as_bigtech_course/week_6/jwt/internal/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

func GenerateToken(info model.UserInfo, secretKey []byte, duration time.Duration) (string, error) {
	claims := model.UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
		Username: info.Username,
		Role:     info.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}

func VerifyToken(tokenStr string, secretKey []byte) (*model.UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&model.UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.New("invalid signing method")
			}

			return secretKey, nil
		},
	)

	if err != nil {
		return nil, errors.Errorf("failed to parse token: %v", err)
	}

	claims, ok := token.Claims.(*model.UserClaims)
	if !ok {
		return nil, errors.Errorf("invalid claims: %v", err)
	}

	return claims, nil
}
