package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/kataras/jwt"
)

type jwtclaim struct {
	Token string `json:"token"`
	Exp   int64  `json:"exp"`
}

func GenerateToken(mySigningKey string) ([]byte, string, error) {
	signingKey := []byte(mySigningKey)
	newToken := uuid.New().String()
	userClaims := jwtclaim{
		Token: newToken,
		Exp:   time.Now().Add(30).Unix(),
	}

	token, err := jwt.Sign(jwt.HS256, signingKey, userClaims)

	return token, newToken, err

}

//ValidateToken validates the jwt token
func ValidateToken(mysecret string, mytoken string) (string, bool) {
	secret := []byte(mysecret)

	tokenString := []byte(mytoken)

	verifiedToken, err := jwt.Verify(jwt.HS256, secret, tokenString)
	if err != nil {
		return "", false
	}
	var claims = struct {
		Token string `json:"token"`
		Exp   int64  `json:"exp"`
	}{}

	verifiedToken.Claims(&claims)
	if claims.Exp < time.Now().Local().Unix() {
		err = errors.New("JWT is expired")
		println(err)
		return claims.Token, false
	}
	// println(claims.Exp) // debug
	// println(claims.Token)

	return "", true

}
