package auth

import (
	"errors"
	"yvpn_server/db"
)

func CreateAccount(u NewUser) (db.Account, error) {
	if err := u.validate(); err != nil {
		return db.Account{}, err
	}

	// Make sure user doesn't already exist
	_, err := db.GetAccountByEmail(u.Email)
	if err == nil {
		return db.Account{}, errors.New("account exists for this email")
	}

	// Hash the password

	// Create the DB record

	// Save the DB record

	// Read the account from the DB

	// Return the account
}
