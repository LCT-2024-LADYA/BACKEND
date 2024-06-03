package utils

import (
	"BACKEND/pkg/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"time"
)

const (
	User    = "user"
	Trainer = "trainer"
	Admin   = "admin"
)

type JWT interface {
	GetExpireTime() time.Duration
	CreateToken(id int, userType string) string
	Authorize(tokenString string, access string) (UserClaim, bool, error)
}

type JWTUtil struct {
	expireTimeOut time.Duration
	secret        string
}

func InitJWTUtil() JWT {
	return JWTUtil{
		expireTimeOut: time.Duration(viper.GetInt(config.JWTExpirationTime)) * time.Hour * 24,
		secret:        viper.GetString(config.JWTSecret),
	}
}

type UserClaim struct {
	jwt.RegisteredClaims
	ID       int
	UserType string
}

func (j JWTUtil) GetExpireTime() time.Duration {
	return j.expireTimeOut
}

func (j JWTUtil) CreateToken(id int, userType string) string {

	expiredAt := time.Now().Add(j.expireTimeOut)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: expiredAt,
			},
		},
		ID:       id,
		UserType: userType,
	})

	signedString, _ := token.SignedString([]byte(j.secret))

	return signedString
}

func (j JWTUtil) Authorize(tokenString string, access string) (UserClaim, bool, error) {
	var claim UserClaim

	token, err := jwt.ParseWithClaims(tokenString, &claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
	if err != nil {
		return UserClaim{}, false, err
	}

	if !token.Valid {
		return UserClaim{}, false, nil
	}

	switch access {
	case User:
		return claim, claim.UserType == User, nil
	default:
		panic("you are passing wrong access")
	}
}
