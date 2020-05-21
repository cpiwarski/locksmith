package ctl

import (
    "bufio"
	"fmt"
    "os"

    "github.com/the-maldridge/locksmith/internal/models"

	"github.com/spf13/cobra"
)

const (
    confPath = "~/.locksmith.conf"
)

var setupCmd = &cobra.Command{
	Use:     "setup",
	Short:   "Create a locksmith profile",
	Long:    `Interactively generate a locksmith profile configuration`,
	Example: "telephone profile setup",
	Run: func(cmd *cobra.Command, args []string) {
        createProfile()
	},
}

func init() {
	profileCmd.AddCommand(setupCmd)
}

// createProfile prompts the user through a series of questions in order to 
// assemble the information needed to write a locksmith profile
func createProfile() error {
    f, err := os.OpenFile(confPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
    if err != nil {
        fmt.Println(err)
        return err
    }
    defer f.Close()

    profile := models.NewProfile()
    scanner := bufio.NewScanner(os.Stdin)

    // collect information
    fmt.Print("Enter a name for the profile: ")
    scanner.Scan()
    profile.Name = scanner.Text()

    fmt.Print("Enter the locksmith server IP: ")
    scanner.Scan()
    profile.Server = scanner.Text()
    return nil
}
