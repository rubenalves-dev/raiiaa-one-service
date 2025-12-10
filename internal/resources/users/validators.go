package users

import (
	"errors"
	"regexp"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

const (
	userFullnameMinLength = 3
	userFullnameMaxLength = 30
)

var (
	ErrFullnameInvalidLength     = errors.New("user fullname is too short")
	ErrFullnameInvalidCharacters = errors.New("user fullname contains invalid characters")
)

type userFullname string

func UserFullname(str string) (userFullname, error) {
	l := utf8.RuneCountInString(str)
	if l < userFullnameMinLength || l > userFullnameMaxLength {
		return "", ErrFullnameInvalidLength
	}
	rgx := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !rgx.MatchString(str) {
		return "", ErrFullnameInvalidCharacters
	}
	return userFullname(str), nil
}

type password string

func Password(str string) (password, error) {
	cost := bcrypt.DefaultCost
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(str), cost)
	if err != nil {
		return "", err
	}
	return password(hashedBytes), nil
}

func VerifyPassword(x string, y string) error {
	return bcrypt.CompareHashAndPassword([]byte(x), []byte(y))
}
