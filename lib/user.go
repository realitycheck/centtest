package centtest

import (
	"github.com/dgrijalva/jwt-go"
)

var (
	jwtAlgorithm = jwt.SigningMethodHS256
	jwtKey       = []byte("secret")
)

type User struct {
	ID string
}

func NewUser(g IDGenerator) *User {
	u := &User{
		ID: g.ID(),
	}
	return u
}

func (u *User) String() string {
	return u.ID
}

func (u *User) token(exp int64) string {
	claims := jwt.MapClaims{
		"sub": u.ID,
		"info": map[string]string{
			"name": u.ID,
		},
	}
	if exp > 0 {
		claims["exp"] = exp
	}
	t, err := jwt.NewWithClaims(jwtAlgorithm, claims).SignedString(jwtKey)
	if err != nil {
		panic(err)
	}
	return t
}
