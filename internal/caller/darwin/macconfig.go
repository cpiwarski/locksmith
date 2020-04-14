package caller

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"github.com/satori/go.uuid"
	"howett.net/plist"
)

var (
	tempFile           = "./config_temp.mobileConfig"
	filename           = "./wireguard.mobileConfig"
	plIDTempl          = "com.%s.wireguard.%s"
	badEncoding        = "&#x[1-9A-Z];"
)

// VendorConfig represents a mobileConfig vendor configuration.
type VendorConfig struct {
	WgQuickConfig string
}

// VpnConfig represents a mobileconfig VPN configuration.
type VpnConfig struct {
	RemoteAddress        string
	AuthenticationMethod string
}

// PayloadContent represents a mobileConfig payload.
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

// MobileConfig represents a WireGuard tunnel mobile configuration.
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
	Config     []byte
	data       MobileConfig
	vpn        VpnConfig
	vendor     VendorConfig
	payload    PayloadContent
	org        string
	tunnelName string
}

// NewConfig creates a new WireGuard mobile configuration.
func NewConfig() *MacConfig {
	return new(MacConfig)
}

// SetOrg sets the name of the MacConfig organization.
func (m *MacConfig) SetOrg(orgName string) {
	m.org = orgName
}

// UpdateVendorConfig resets the information in the VendorConfig.
func (m *MacConfig) UpdateVendorConfig() {
	vc := VendorConfig{
		WgQuickConfig: string(m.Config),
	}
	m.vendor = vc
}

// UpdateVpnConfig updates the vpn config for the MacConfig.
func (m *MacConfig) UpdateVpnConfig(newAddr string) {
	vc := VpnConfig{
		RemoteAddress:        newAddr,
		AuthenticationMethod: "Password",
	}
	m.vpn = vc
}

// UpdatePayload updates the payload data for the struture
func (m *MacConfig) UpdatePayload() {
	i := fmt.Sprintf(plIDTempl, m.org, m.tunnelName)
	id := uuid.NewV4() 
	p := PayloadContent{
		PayloadDisplayName: "VPN",
		PayloadType:        "com.apple.vpn.managed",
		PayloadVersion:     1,
		PayloadIdentifier:  i,
		PayloadUUID:        id.String(),
		UserDefinedName:    m.tunnelName,
		VPNType:            "VPN",
		VPNSubType:         "com.wireguard.macos",
		VendorConfig:       m.vendor,
		VPN:                m.vpn,
	}
	m.payload = p
}

// UpdateMobileConfig updates the mobileConfig plist stored in the
// MacConfig.
func (m *MacConfig) UpdateMobileConfig() {
	id1 := uuid.NewV4()
	id2 := uuid.NewV4()
	cIdentifier := fmt.Sprintf(plIDTempl, m.org, id1.String())
	c := MobileConfig{
		PayloadDisplayName: "Tunnel",
		PayloadType:        "Configuration",
		PayloadVersion:     1,
		PayloadIdentifier:  cIdentifier,
		PayloadUUID:        id2.String(),
		PayloadContent:     m.payload,
	}
	m.data = c
}

// InstallConfig installs the tunnel configuration onto the Mac.
func (m *MacConfig) InstallConfig() error {
	
	plist, err := plist.MarshalIndent(m.data, plist.XMLFormat, "\t")
	if err != nil {
		return err
	}
	fmt.Print(string(plist))

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
	nf, err := os.OpenFile(outFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	defer nf.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		newlineByte := []byte("\n")
		nextLine := []byte(s.Text())
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
