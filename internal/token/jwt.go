package token

import (
	"fmt"
	"strconv"
	"time"
	"github.com/AyuorusAguilar/Biscord/internal/state"
	"github.com/golang-jwt/jwt/v5"
)

func CreateJWT(userID int, s state.State) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.RegisteredClaims{
		Issuer: "Biscord Handsome Staff",
		Subject: strconv.Itoa(userID),
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
	})
	tokenstring, err := token.SignedString([]byte(s.K))
	if err != nil {
		return "", fmt.Errorf("Error creating token")
	}
	return tokenstring, nil
}
func ValidateJWT(tokenString string, tokenSecret string) (int, error) {
	a := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &a,
		func(t *jwt.Token) (any, error) { return []byte(tokenSecret), nil }
	)
	if err != nil {
		return 0, err
	}
	subject, err := token.Claims.GetSubject()
	if err != nil {
		return 0, err
	}
	id, err := strconv.Atoi(subject)
	if err != nil {
		return 0, err
	}
	return id, nil
}