package tests

import (
	"fmt"
	"github.com/Phanile/go-exchange-auth/tests/suite"
	authv1 "github.com/Phanile/go-exchange-protos/generated/go/auth"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

const defaultPasswordLen = 8

func TestAuth_Login(t *testing.T) {
	ctx, st := suite.NewSuite(t)

	email := gofakeit.Email()
	password := generatePassword()

	respReg, errReg := st.AuthClient.Register(ctx, &authv1.RegisterRequest{
		Email:    email,
		Password: password,
	})

	require.NoError(t, errReg)
	assert.NotEmpty(t, respReg.GetUserId())

	respLogin, errLogin := st.AuthClient.Login(ctx, &authv1.LoginRequest{
		Email:    email,
		Password: password,
	})

	require.NoError(t, errLogin)

	privatePem := os.Getenv("JWT_PRIVATE_KEY")
	require.NotEmpty(t, privatePem)

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privatePem))
	require.NoError(t, err)

	tokenParsed, errParse := jwt.Parse(respLogin.GetToken(), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return &privateKey.PublicKey, nil
	})

	require.NoError(t, errParse)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, respReg.GetUserId(), int64(claims["uid"].(float64)))
}

func generatePassword() string {
	return gofakeit.Password(true, true, true, true, false, defaultPasswordLen)
}
