package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserId string `json:"user_id"`
	jwt.RegisteredClaims
}

func CreateToken(userId string, secretKey []byte) (string, error) {
	claims := Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24))
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string, secretKey []byte) (*Claims, error){
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil{
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid{
		return claims, nil
	} else {  
		return nil, err
	}
}