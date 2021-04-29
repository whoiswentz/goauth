package helpers

import (
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type SignedDetails struct {
	Email string
	Name  string
	Id    int64
	jwt.StandardClaims
}

var (
	SECRET_KEY      = RandomStringBase64(256)
	ErrTokenExired  = errors.New("token is expired")
	ErrTokenInvalid = errors.New("the token is invalid")
)

func GenerateAllTokens(email string, name string, id int64) (signedToken string, signedRefreshToken string, err error) {
	TOKEN_DURATION := os.Getenv("TOKEN_DURATION")
	REFRESH_DURATION := os.Getenv("REFRESH_DURATION")

	duration, err := strconv.Atoi(TOKEN_DURATION)
	refreshDuration, err := strconv.Atoi(REFRESH_DURATION)
	if err != nil {
		return "", "", err
	}

	claims := &SignedDetails{
		Email: email,
		Name:  name,
		Id:    id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(duration)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(refreshDuration)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken, nil
}

func ValidateToken(signedToken string) (*SignedDetails, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		return nil, ErrTokenInvalid
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, ErrTokenExired
	}

	return claims, nil
}
