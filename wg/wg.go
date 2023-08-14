package wg

import (
	"fmt"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"strings"
)

type Keys struct {
	Public  wgtypes.Key
	Private wgtypes.Key
}

// GenerateKeys returns a private key and its corresponding public key.
func GenerateKeys() (wgtypes.Key, wgtypes.Key, error) {
	privateKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return wgtypes.Key{}, wgtypes.Key{}, err
	}

	publicKey := privateKey.PublicKey()

	return privateKey, publicKey, nil
}

func GenerateServerConfig(servKeys Keys, clients map[string]Keys) (string, error) {
	// Begin with the server's [Interface] configuration
	configBuilder := strings.Builder{}

	configBuilder.WriteString("[Interface]\n")
	configBuilder.WriteString(fmt.Sprintf("Address = 10.0.0.1/24\n"))
	configBuilder.WriteString(fmt.Sprintf("ListenPort = 51820\n"))
	configBuilder.WriteString(fmt.Sprintf("PrivateKey = %s\n", servKeys.Private))
	configBuilder.WriteString("PostUp = iptables -A FORWARD -i %i -j ACCEPT; iptables -A FORWARD -o %i -j ACCEPT; iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE\n")
	configBuilder.WriteString("PostDown = iptables -D FORWARD -i %i -j ACCEPT; iptables -D FORWARD -o %i -j ACCEPT; iptables -t nat -D POSTROUTING -o eth0 -j MASQUERADE\n")
	configBuilder.WriteString("DNS = 10.0.0.1\n\n")

	// For each client, append a [Peer] configuration
	for allowedIP, keys := range clients {
		configBuilder.WriteString("[Peer]\n")
		configBuilder.WriteString(fmt.Sprintf("PublicKey = %s\n", keys.Public))
		configBuilder.WriteString(fmt.Sprintf("AllowedIPs = %s\n\n", allowedIP))
	}

	return configBuilder.String(), nil
}

func GenerateCloudInit(wgConfig string) string {
	const cloudInitTmpl = `#cloud-config
packages:
  - wireguard
  - iptables-persistent

write_files:
  - content: |
      %s
    path: /etc/wireguard/wg0.conf
    permissions: '0600'

runcmd:
  - wg-quick up wg0
  - systemctl enable wg-quick@wg0.service
  - iptables -A INPUT -p udp -m udp --dport 51820 -j ACCEPT
  - iptables -A FORWARD -i wg0 -j ACCEPT
  - iptables-save > /etc/iptables/rules.v4

`

	return fmt.Sprintf(cloudInitTmpl, wgConfig)
}
