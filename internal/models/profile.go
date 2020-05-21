package models

type Profile struct {
	name       string
	uPrivKey   string
	uPublicKey string
	PrivateKey string
	serverAddr string
	Address    string
	Pubkey     string
	AllowedIPs string
	Endpoint   string
	tunnel     InterfaceConfig
}
