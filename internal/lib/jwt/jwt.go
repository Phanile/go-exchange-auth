package jwt

import (
	"github.com/Phanile/go-exchange-auth/internal/domain/models"
	"github.com/golang-jwt/jwt"
	"time"
)

func NewToken(user *models.User, duration time.Duration, privateKey string) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["uid"] = user.Id
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	signedToken, err := token.SignedString([]byte(privateKey))

	if err != nil {
		return ""
	}

	return signedToken
}
