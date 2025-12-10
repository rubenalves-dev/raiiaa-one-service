package validators

import (
	"errors"
	"net/mail"
)

var (
	ErrInvalidEmail = errors.New("invalid email address")
)

type email string

func Email(str string) (email, error) {
	_, err := mail.ParseAddress(str)
	if err != nil {
		return "", ErrInvalidEmail
	}
	return email(str), nil
}
