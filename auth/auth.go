package auth

import (
	"errors"
	"regexp"
)

type NewUser struct {
	Email       string
	Password    string
	ConfirmPass string
	InviteCode  string
}

func (u *NewUser) validate() error {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !re.MatchString(u.Email) {
		return errors.New("invalid email address")
	}

	re = regexp.MustCompile(`(?=.*\d)(?=.*[a-z])(?=.*[A-Z]).{8,}`)
	if !re.MatchString(u.Password) {
		return errors.New("password must contain at least one number, one uppercase and lowercase letter, and at least 8 characters")
	}

	if u.Password != u.ConfirmPass {
		return errors.New("password and confirm password do not match")
	}

	return nil
}
