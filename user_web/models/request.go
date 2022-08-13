package models

import "github.com/dgrijalva/jwt-go"

type CustomClaims struct {
	ID          uint
	NickName    string
	AuthorityId uint // 对应数据库中的 role
	jwt.StandardClaims
}
