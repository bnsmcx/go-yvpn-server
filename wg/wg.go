package wg

import "golang.zx2c4.com/wireguard/wgctrl/wgtypes"

// GenerateKeys returns a private key and its corresponding public key.
func GenerateKeys() (wgtypes.Key, wgtypes.Key, error) {
	privateKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return wgtypes.Key{}, wgtypes.Key{}, err
	}

	publicKey := privateKey.PublicKey()

	return privateKey, publicKey, nil
}
