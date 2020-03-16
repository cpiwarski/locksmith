package caller

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"howett.net/plist"
	"github.com/satori/go.uuid"
)

var (
	tempFile             = "./config_temp.mobileConfig"
	filename             = "./wireguard.mobileConfig"
	identifierTemplate   = "com.%s.wireguard.locksmith.%s"
	plIdTemp = "com.%s.wireguard.%s"
	badEncoding          = "&#x[1-9A-Z];"
)

// vendorConfig represents a mobileConfig vendor configuration.
type VendorConfig struct {
	WgQuickConfig string
}

// vpnConfig represents a mobileconfig VPN configuration.
type VpnConfig struct {
	RemoteAddress        string
	AuthenticationMethod string
}

// payloadContent represents a mobileConfig payload.
type PayloadContent struct {
	PayloadDisplayName string
	PayloadType        string
	PayloadVersion     int
	PayloadIdentifier  string
	PayloadUUID        string
	UserDefinedName    string
	VPNType            string
	VPNSubType         string
	VendorConfig       VendorConfig
	VPN                VpnConfig
}

// mobileConfig represents a WireGuard tunnel mobile configuration.
type MobileConfig struct {
	PayloadDisplayName string
	PayloadType        string
	PayloadVersion     int
	PayloadIdentifier  string
	PayloadUUID        string
	PayloadContent     PayloadContent
}

// MacConfig represents a MacOS X WireGuard configuration.
type MacConfig struct {
	Config      []byte
	data        MobileConfig
	vpn         VpnConfig
	vendor      VendorConfig
	payload     PayloadContent
	config      MobileConfig
	org         string
	tunnelName string
}

// NewConfig creates a new WireGuard mobile configuration.
func NewConfig() *MacConfig {
	return new(MacConfig)
}

// RecompileMobileConfig inserts all the data into the MacConfig data 
// structures
func (m *MacConfig) RecompileMobileConfig() {

	configTemplate := MobileConfig {
		PayloadDisplayName: "Tunnel",
		PayloadType:        "Configuration",
		PayloadVersion:     1,
		PayloadIdentifier:  "",
		PayloadUUID:        "",
		PayloadContent:     defaultPayload,
	}
	m.config = configTemplate
	wgTunnelConfig, err := plist.Marshal(configTemplate, plist.XMLFormat)
	if err != nil {
		return MacConfig{}
	}

	newMacConfig := MacConfig {
		localDevice: newDevice,
		tunnel:      wgTunnelConfig,
		data:        configTemplate,
	}
}

// SetOrg sets the name of the MacConfig organization.
func (m *MacConfig) SetOrg(orgName string) {
	m.org = orgName
}

// UpdateVendorConfig resets the information in the VendorConfig.
func (m *MacConfig) UpdateVendorConfig() {
	vc := VendorConfig {
		WgQuickConfig: string(m.Config),
	}
	m.vendor = vc
}

// UpdateVpnConfig updates the vpn config for the MacConfig.
func (m *MacConfig) UpdateVpnConfig() {
	vc := VpnConfig{
		RemoteAddress:        newAddr,
		AuthenticationMethod: "Password",
	}
	m.vpn = vc
}

// UpdatePayload updates the payload data for the struture
func (m *MacConfig) UpdatePayload() error {
	i := fmt.Sprintf(plIdTempl, m.org, m.tunnelName)
	id, err = CreateUUID()
	if err != nil {
		return err
	}
	p := PayloadContent {
		PayloadDisplayName: "VPN",
		PayloadType:        "com.apple.vpn.managed",
		PayloadVersion:     1,
		PayloadIdentifier:  i,
		PayloadUUID:        id,
		UserDefinedName:    m.tunnelName,
		VPNType:            "VPN",
		VPNSubType:         "com.wireguard.macos",
		VendorConfig:       m.vendor,
		VPN:                m.vpn,
	}
	m.payload = p
	return nil
}

// UpdateMobileConfig updates the mobileConfig plist stored in the
// MacConfig.
func (m *MacConfig) UpdateMobileConfig() error {
	id1, err := CreateUUID()
	if err != nil {
		return err
	}
	id2, err := CreateUUID()
	if err != nil {
		return err
	}
	configIdentifier := fmt.Sprintf(plIdentifierTemplate, m.org, id1)
	c := MobileConfig{
		PayloadDisplayName: "Tunnel",
		PayloadType:        "Configuration",
		PayloadVersion:     1,
		PayloadIdentifier:  configIdentifier,
		PayloadUUID:        id2,
		PayloadContent:     m.payload,
	}
	m.data = c
	return nil
}

// InstallConfig installs the tunnel configuration onto the Mac.
func (m *MacConfig) InstallConfig() error {

	// TODO: Install Config without shell commands if possible

	return nil
}

// UninstallConfig removes the tunnel configuration from the Mac.
func (m *MacConfig) UninstallConfig() error {

	// TODO: Remove installation without using shell commands if possible
	return nil
}

// WriteMobileConfig  writes a new configuration file for installation.
func WriteMobileConfig(mobileConfig []uint8) (string, error) {

	_ = os.Remove(tempFile)

	f, err := os.OpenFile(tempFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0600)
	if err != nil {
		return "", err
	}

	defer f.Close()

	_, err = f.Write(mobileConfig)
	if err != nil {
		return "", err
	}
	return filename, nil
}

// removeBadCharacters removes any newline characters that rendered poorly
// in the mobileConfig.
func removeBadCharacters(inFile, outFile string) error {
	garbage := regexp.MustCompile(badEncoding)

	_ = os.Remove(outFile)

	// Open temporary file
	f, err := os.Open(inFile)
	if err != nil {
		return err
	}

	defer f.Close()

	// Open final config file
	nf, err := os.OpenFile(outFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0600)
	if err != nil {
		return err
	}

	defer nf.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		newlineByte := []byte("\n")
		nextLine := []byte(scanner.Text())
		regScan := garbage.ReplaceAll(nextLine, newlineByte)
		_, err = nf.Write(regScan)

		if err != nil {
			return err
		}
	}
	_ = os.Remove(inFile)
	return nil
}

// RemoveLocalFiles removes the local files written during the config
// installation process.
func RemoveLocalFiles() {
	_ = os.Remove(tempFile)
	_ = os.Remove(filename)
}

// CreateUUID returns a string UUID or an error.
func CreateUUID() (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return id, nil
}
