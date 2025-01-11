package helper

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func GetBcryptPassword(password string) (BcryptPassword string, err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func AnalysisBcryptPassword(BcryptPassword, password string) bool {
	fmt.Println("校验密码")
	fmt.Println(bcrypt.CompareHashAndPassword([]byte(BcryptPassword), []byte(password)) == nil)
	return bcrypt.CompareHashAndPassword([]byte(BcryptPassword), []byte(password)) == nil
}
