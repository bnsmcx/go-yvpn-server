package auth

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"yvpn_server/db"
)

func CreateAccount(u NewUser) (*db.Account, error) {
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
	if err == nil {
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
