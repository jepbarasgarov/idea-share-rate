package helpers

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
)

type AuthTokenClaim struct {
	Username string
	jwt.StandardClaims
}

func GenerateAccessToken(username string, duration int, secret string) (accessToken string, err error) {
	clog := log.WithFields(log.Fields{
		"method": "GenerateAccessToken",
	})

	expiresAt := time.Now().Add(time.Second * time.Duration(duration)).Unix()

	claims := &AuthTokenClaim{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, err = token.SignedString([]byte(secret))
	if err != nil {
		eMsg := "error in token.SignedString"
		clog.WithError(err).Error(eMsg)
		return
	}

	return
}

func VerifyAccessToken(token string, secret string) (username string, err error) {

	clog := log.WithFields(log.Fields{
		"method": "VerifyAccessToken",
	})

	claims := jwt.MapClaims{}

	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	username = fmt.Sprintf("%v", claims["Username"])

	if err != nil {
		eMsg := "error in jwt.ParseWithClaims"
		clog.WithError(err).Error(eMsg)
		return username, err
	}

	return username, err
}

func GetUsernameFromAccessToken(token string, secret string) (username string, err error) {

	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	username = fmt.Sprintf("%v", claims["Username"])

	if err != nil {
		return
	}

	return
}

func GenerateRefreshToken() (refreshToken string, err error) {
	x, err := uuid.NewV4()
	if err != nil {
		fmt.Println("error RefreshToken")
		return
	}

	refreshToken = x.String()

	return
}
