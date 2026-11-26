package token

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CreateJWT creates a JWT xdxd (Javascript Web Token or
// something like that). Use the JWT duration especified
// in the state scruct unless you have good reason not to.
func CreateJWT(userID int, tokenSecret string, expireIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: "Biscord Handsome Staff",
		Subject: strconv.Itoa(userID),
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireIn)),
	})
	tokenstring, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("Error creating token")
	}
	return tokenstring, nil
}

// ValidateJWT returns the id of the subject and a posible error.
// Notice that if the error includes 'token is expired' you should
// return a *401 Unauthorized* response code instead of anything else,
// so the client can request a new jwt and try the request again.
func ValidateJWT(tokenString string, tokenSecret string) (int64, error) {
	a := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &a,
		func(t *jwt.Token) (any, error) { return []byte(tokenSecret), nil },
	)
	if err != nil {
		return 0, err
	}
	subject, err := token.Claims.GetSubject()
	if err != nil {
		return 0, err
	}
	id, err := strconv.ParseInt(subject, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// CreateRefreshToken returns a the raw token and a hash.
// Both in a 32 bytes format.
func CreateRefreshToken() ([]byte, [32]byte, error) {
	var raw = make([]byte, 32)
	_, err := rand.Read(raw)
	if err != nil {
		return make([]byte, 0), [32]byte{}, fmt.Errorf("Error trying to read random data: %v", err)
	}
	hash := sha256.Sum256(raw)
	return raw, hash, nil
}

func GetHash(s []byte) ([32]byte) {
	return sha256.Sum256(s)
}