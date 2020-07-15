package common

import (
	"com.jumbo/ginessential/model"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtKey = []byte("a_secret_crect")

type Claims struct {
	UserId uint
	jwt.StandardClaims
}

func ReleaseTken(user model.User) (string, error)  {
	expirationTime := time.Now().Add(7 *24 *time.Hour)
	claims := Claims{
		UserId:         user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt: time.Now().Unix(),
			Issuer: "developer.jumbo",
			Subject: "user token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	return token, claims, err
}

/*
token: 三部分组成由.分隔 ` echo eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9 | base64 -d` 查看
第一部分储存加密协议
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.  -> {"alg":"HS256","typ":"JWT"}%
第二部分储存jwtClaims中信息payload
eyJVc2VySWQiOjQsImV4cCI6MTU5NTM4NjkwOCwiaWF0IjoxNTk0NzgyMTA4LCJpc3MiOiJkZXZlbG9wZXIuanVtYm8iLCJzdWIiOiJ1c2VyIHRva2VuIn0.
-> {"UserId":4,"exp":1595386908,"iat":1594782108,"iss":"developer.jumbo","sub":"user token%
第三部分由前两部分加key哈希的值
oj4K24zhzAxIgEoQhOgvZBSYJ0wKhvf-kQ4fk_IFx6E
*/