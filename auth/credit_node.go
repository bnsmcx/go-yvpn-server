package auth

import (
	"errors"
	"github.com/google/uuid"
	"os"
	"yvpn_server/db"
)

type NewCreditNode struct {
	InviteCode string
}

func (n *NewCreditNode) Create() (*db.Account, error) {
	// Validate Invite Code
	if n.InviteCode != os.Getenv("YVPN_INVITE_CODE") {
		return nil, errors.New("invalid invite code")
	}

	// Create and save the DB record
	newRecord := db.Account{
		ID: uuid.New(),
	}

	err := newRecord.Save()
	if err != nil {
		return &db.Account{}, err
	}

	// Read the account from the DB
	account, err := db.GetAccount(newRecord.ID)
	if err != nil {
		return &db.Account{}, errors.New("error retrieving account")
	}

	return account, nil
}
