package auth

import (
	"errors"
	"github.com/google/uuid"
	"os"
	"yvpn_server/db"
)

type NewCreditNode struct {
	InviteCode        string
	DigitalOceanToken string
}

func (n *NewCreditNode) Create() (*Account, error) {
	// Validate Invite Code
	if n.InviteCode != os.Getenv("YVPN_INVITE_CODE") {
		return nil, errors.New("invalid invite code")
	}

	// Create and save the DB record
	newRecord := db.Account{
		ID: uuid.New(),
	}

	if err := newRecord.Save(); err != nil {
		return nil, err
	}

	a := Account{
		DigitalOceanToken: n.DigitalOceanToken,
		ID:                newRecord.ID,
	}

	return &a, nil
}
