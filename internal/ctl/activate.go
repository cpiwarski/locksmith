package ctl

import (
	"fmt"
	"os"
	"text/template"

	"github.com/spf13/cobra"
)

const (
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

var activateCmd = &cobra.Command {
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

// activate executes all tasks associated with activating a locksmith profile
func activate() {

	f, err := os.Create(configPath)
	if err != nil {
		fmt.Println(err)
	}
    defer f.Close()
}
