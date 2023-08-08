package auth

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"yvpn_server/db"
)

type NewAccount struct {
	Email       string
	Password    string
	ConfirmPass string
	InviteCode  string
}

func (u *NewAccount) validate() error {
	// Email validation
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !re.MatchString(u.Email) {
		return errors.New("invalid email address")
	}

	// Password validation for complexity
	if len(u.Password) < 8 ||
		!regexp.MustCompile(`[a-z]`).MatchString(u.Password) ||
		!regexp.MustCompile(`[A-Z]`).MatchString(u.Password) ||
		!regexp.MustCompile(`\d`).MatchString(u.Password) {
		return errors.New("password must contain at least one number, one uppercase and lowercase letter, and at least 8 characters")
	}

	// Confirm password validation
	if u.Password != u.ConfirmPass {
		return errors.New("password and confirm password do not match")
	}

	return nil
}

func (u *NewAccount) Create() (*db.Account, error) {
	if err := u.validate(); err != nil {
		return &db.Account{}, err
	}

	// Make sure user doesn't already exist
	_, err := db.GetAccountByEmail(u.Email)
	if err == nil {
		return &db.Account{}, errors.New("account exists for this email")
	}

	// Hash the password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return &db.Account{}, err
	}

	// Create and save the DB record
	newRecord := db.Account{
		ID:       uuid.New(),
		Email:    u.Email,
		Password: hashedPass,
	}

	err = newRecord.Save()
	if err != nil {
		return &db.Account{}, err
	}

	// Read the account from the DB
	account, err := db.GetAccountByEmail(u.Email)
	if err != nil {
		return &db.Account{}, errors.New("error retrieving account")
	}

	return account, nil
}
