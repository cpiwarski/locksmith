package ctl

import (
	"fmt"
	"os"
	"text/template"

	"github.com/the-maldridge/locksmith/internal/models"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	installPath = "~/.locksmith"
	configName  = ".env"
	configPath  = "."
	configType  = "yaml"
)

var (
	tmpl   = template.Must(template.New("Tunnel").Parse(tunnel))
	tunnel = `[Interface]
PrivateKey = {{.PrivateKey}}
Address = {{.Address}}

[Peer]
PublicKey = {{.Pubkey}}
AllowedIPs = {{.AllowedIPs}}
Endpoint = {{.Endpoint}}
`
)

var activateCmd = &cobra.Command{
	Use:     "activate",
	Short:   "Enable a locksmith profile",
	Long:    `Enable the supplied locksmith profile.`,
	Example: "telephone activate <profile>",
	Run: func(cmd *cobra.Command, args []string) {
		activate()
	},
}

func init() {
	rootCmd.AddCommand(activateCmd)
}

func viperInit() {
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath(configPath)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Viper Config Failure")
		fmt.Println(err)
		return
	}
	return
}

func tunnelInit() *models.Profile {
	profile := &models.Profile{
		PrivateKey: viper.GetString("uPrivateKey"),
		Address:    viper.GetString("address"),
		Pubkey:     viper.GetString("publicKey"),
		AllowedIPs: viper.GetString("allowedIPs"),
		Endpoint:   viper.GetString("endpoint"),
	}
	return profile
}

// activate executes all tasks associated with activating a locksmith profile
func activate() {
	viperInit()
	tunnelSt := tunnelInit()

	f, err := os.Create(installPath)
	if err != nil {
		fmt.Println(err)
	}

	_ = tmpl.Execute(f, tunnelSt)
}
