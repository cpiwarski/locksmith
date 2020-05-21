package ctl

import (
    "bufio"
	"fmt"
    "io"
    "os"

	"github.com/spf13/cobra"
)

const (
    confPath = "./locksmith.conf"
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

type Profile struct {
    Name string
    Tunnel string
    Server string
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

    profile := Profile{}
    scanner := bufio.NewScanner(os.Stdin)

    // collect information
    fmt.Print("Enter a name for the profile: ")
    scanner.Scan()
    profile.Name = scanner.Text()

    fmt.Print("Enter the locksmith server IP: ")
    scanner.Scan()
    profile.Server = scanner.Text()

    fmt.Print("Enter the locksmith tunnel: ")
    scanner.Scan()
    profile.Tunnel = scanner.Text()

    // write information to conf
    if _, err = io.WriteString(f, profile.Name + "\n"); err != nil {
        return err
    }
    if _, err = io.WriteString(f, profile.Server + "\n"); err != nil {
        return err
    }
    if _, err = io.WriteString(f, profile.Tunnel + "\n"); err != nil {
        return err
    }

    return nil
}
