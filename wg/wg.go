package wg

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/skip2/go-qrcode"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"log"
	"strings"
	"text/template"
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

func GenerateClientConfig(
	serverIP, clientIP string,
	serverPubKey, clientPrivKey string) (string, error) {

	configBuilder := strings.Builder{}

	configBuilder.WriteString("[Interface]\n")
	configBuilder.WriteString(fmt.Sprintf("Address = %s/24\n", clientIP))
	configBuilder.WriteString("DNS = 1.1.1.1\n")
	configBuilder.WriteString(fmt.Sprintf("PrivateKey = %s\n\n", clientPrivKey))

	configBuilder.WriteString("[Peer]\n")
	configBuilder.WriteString(fmt.Sprintf("PublicKey = %s\n", serverPubKey))
	configBuilder.WriteString(fmt.Sprintf("Endpoint = %s:51820\n", serverIP))
	configBuilder.WriteString("AllowedIPs = 0.0.0.0/0\n")

	return configBuilder.String(), nil
}

func GenerateServerConfig(servKeys Keys, clients map[string]Keys) (string, error) {
	// Begin with the server's [Interface] configuration
	configBuilder := strings.Builder{}

	configBuilder.WriteString("[Interface]\n")
	configBuilder.WriteString(fmt.Sprintf("Address = 10.0.0.1/24\n"))
	configBuilder.WriteString(fmt.Sprintf("ListenPort = 51820\n"))
	configBuilder.WriteString(fmt.Sprintf("PrivateKey = %s\n", servKeys.Private))
	configBuilder.WriteString("PostUp = iptables -A FORWARD -i %i -j ACCEPT; iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE\n")
	configBuilder.WriteString("PostDown = iptables -D FORWARD -i %i -j ACCEPT; iptables -t nat -D POSTROUTING -o eth0 -j MASQUERADE\n\n")

	// For each client, append a [Peer] configuration
	for i := 2; i <= 255; i++ {
		allowedIP := fmt.Sprintf("10.0.0.%d", i)
		keys, ok := clients[allowedIP]
		if !ok {
			log.Println("missing config for: ", allowedIP)
			continue
		}
		configBuilder.WriteString("[Peer]\n")
		configBuilder.WriteString(fmt.Sprintf("PublicKey = %s\n", keys.Public))
		configBuilder.WriteString(fmt.Sprintf("AllowedIPs = %s\n\n", allowedIP))
	}

	return configBuilder.String(), nil
}
func GenerateCloudInit(wgConfig string) (string, error) {
	const cloudInitTmpl = `#cloud-config

packages:
  - iptables-persistent

write_files:
  - content: |
{{ indent 6 .ConfigContent }}
    path: /etc/wireguard/wg0.conf
    permissions: '0600'

runcmd:
  - add-apt-repository -y ppa:wireguard/wireguard
  - apt update
  - apt install -y wireguard-tools
  - apt install -y resolvconf
  - sysctl -w net.ipv4.ip_forward=1
  - wg-quick up wg0
  - systemctl enable wg-quick@wg0.service
`

	data := map[string]string{
		"ConfigContent": wgConfig,
	}

	tmpl, err := template.New("cloud-init").Funcs(template.FuncMap{
		"indent": func(i int, input string) string {
			prefix := strings.Repeat(" ", i)
			return prefix + strings.ReplaceAll(input, "\n", "\n"+prefix)
		},
	}).Parse(cloudInitTmpl)

	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func GenerateQR(config string) (string, error) {
	// Generate QR code
	qrCode, err := qrcode.New(config, qrcode.Medium)
	if err != nil {
		return "", err
	}

	qr, err := qrCode.PNG(300)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(qr), nil
}
