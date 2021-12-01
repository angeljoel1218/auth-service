package fixture

import (
	"auth-service/config"
	apierror "auth-service/src/domain/apierrors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type TokenCustomClaims struct {
	UserId *int64 `json:"user_id"`
	jwt.StandardClaims
}

func CreateToken(id int64) (string, error) {
	now := time.Now().Local()

	claims := TokenCustomClaims{
		UserId: &id,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(time.Hour * time.Duration(config.C.Token.ExpiredHour)).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(config.C.Token.Secret))

	return signedToken, err
}

func DecodeToken(tokenPart string) (*TokenCustomClaims, error) {
	claims := &TokenCustomClaims{}
	token, err := jwt.ParseWithClaims(tokenPart, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.C.Token.Secret), nil
	})

	if err != nil {
		apier := apierror.NewErrorApi(apierror.Unauthorized, err.Error())
		return nil, apier
	}

	if !token.Valid {
		apier := apierror.NewErrorApi(apierror.Unauthorized, "ID token is invalid")
		return nil, apier
	}

	claims, ok := token.Claims.(*TokenCustomClaims)

	if !ok {
		apier := apierror.NewErrorApi(apierror.Unauthorized, "ID token valid but couldn't parse claims")
		return nil, apier
	}

	return claims, err
}
