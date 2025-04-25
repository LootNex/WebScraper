package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func GetUserID(token string) (string, error) {
	secret := os.Getenv("SECRET")

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, flag := token.Method.(*jwt.SigningMethodHMAC); !flag {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, flag := parsedToken.Claims.(jwt.MapClaims); flag && parsedToken.Valid {
		userID, flag := claims["uuid"].(string)
		if !flag {
			return "", fmt.Errorf("invalid user id")
		}

		if exp, flag := claims["exp"].(float64); flag {
			expirationTime := time.Unix(int64(exp), 0)
			if expirationTime.Before(time.Now()) {
				return "", fmt.Errorf("token has expired")
			}
		} else {
			return "", fmt.Errorf("missing expiration claim")
		}

		return userID, nil
	}

	return "", fmt.Errorf("invalid token")
}
