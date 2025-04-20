package jwt

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/auth/internal/domain/models"
)

func NewToken(user models.User, duration time.Duration) (string, error) {
	secret := os.Getenv("SECRET")

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uuid"] = user.ID
	claims["login"] = user.Login
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
