package caller

type WindowsConfig struct {
	DeviceName string
	org string
	Config []byte
}

// NewWindowsConfig creates a new Windows configuration.
func NewConfig() WindowsConfig {
	return new(WindowsConfig)
}

// InstallConfig installs the configuration to the Windows device.
func (w WindowsConfig) InstallConfig() error {
	//TODO: Complete method
	return nil
}

// UninstallConfig removes the configuration from the Windows device.
func (w WindowsConfig) UninstallConfig() error {
	//TODO: Complete method
	return nil
}
