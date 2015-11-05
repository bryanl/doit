package commands

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"

	"github.com/digitalocean/doctl"
	"github.com/digitalocean/godo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configFile = ".doctlcfg"
)

var (
	// DoctlCmd is the base command.
	DoctlCmd = &cobra.Command{
		Use: "doctl",
	}

	// Token holds the global authorization token.
	Token string

	// Output holds the global output format.
	Output string

	writer = os.Stdout
)

func init() {
	viper.SetConfigType("yaml")

	DoctlCmd.PersistentFlags().StringVarP(&Token, "access-token", "t", "", "DigtialOcean API V2 Access Token")
	DoctlCmd.PersistentFlags().StringVarP(&Output, "output", "o", "text", "output formt [text|json]")
}

// LoadConfig loads out configuration.
func LoadConfig() error {
	fp, err := configFilePath()
	if err != nil {
		return fmt.Errorf("can't find home directory: %v", err)
	}
	if _, err := os.Stat(fp); err == nil {
		file, err := os.Open(fp)
		if err != nil {
			return fmt.Errorf("can't open configuration file %q: %v", fp, err)
		}
		viper.ReadConfig(file)
	}

	return nil
}

// Execute executes the base command.
func Execute() {
	initializeConfig()
	addCommands()
	DoctlCmd.Execute()
}

// AddCommands adds sub commands to the base command.
func addCommands() {
	DoctlCmd.AddCommand(Account())
	DoctlCmd.AddCommand(Actions())
	DoctlCmd.AddCommand(Domain())
	DoctlCmd.AddCommand(DropletAction())
	DoctlCmd.AddCommand(Droplet())
	DoctlCmd.AddCommand(FloatingIP())
	DoctlCmd.AddCommand(FloatingIPAction())
	DoctlCmd.AddCommand(Images())
	DoctlCmd.AddCommand(Region())
	DoctlCmd.AddCommand(Size())
	DoctlCmd.AddCommand(SSHKeys())
	DoctlCmd.AddCommand(SSH())
}

func initFlags() {
	viper.SetEnvPrefix("DIGITALOCEAN")
	viper.BindEnv("access-token", "DIGITALOCEAN_ACCESS_TOKEN")
	viper.BindPFlag("access-token", DoctlCmd.PersistentFlags().Lookup("access-token"))
	viper.BindPFlag("output", DoctlCmd.PersistentFlags().Lookup("output"))
}

func loadDefaultSettings() {
	viper.SetDefault("output", "text")
}

// InitializeConfig initializes the doctl configuration.
func initializeConfig() {
	loadDefaultSettings()
	LoadConfig()
	initFlags()

	if DoctlCmd.PersistentFlags().Lookup("access-token").Changed {
		viper.Set("access-token", Token)
	}

	if DoctlCmd.PersistentFlags().Lookup("output").Changed {
		viper.Set("output", Output)
	}
}

func configFilePath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(usr.HomeDir, configFile)
	return dir, nil
}

func addStringFlag(cmd *cobra.Command, name, dflt, desc string) {
	fn := flagName(cmd, name)
	cmd.Flags().String(name, dflt, desc)
	viper.BindPFlag(fn, cmd.Flags().Lookup(name))
}

func addIntFlag(cmd *cobra.Command, name string, def int, desc string) {
	fn := flagName(cmd, name)
	cmd.Flags().Int(name, def, desc)
	viper.BindPFlag(fn, cmd.Flags().Lookup(name))
}

func addBoolFlag(cmd *cobra.Command, name string, def bool, desc string) {
	fn := flagName(cmd, name)
	cmd.Flags().Bool(name, def, desc)
	viper.BindPFlag(fn, cmd.Flags().Lookup(name))
}

func addStringSliceFlag(cmd *cobra.Command, name string, def []string, desc string) {
	fn := flagName(cmd, name)
	cmd.Flags().StringSlice(name, def, desc)
	viper.BindPFlag(fn, cmd.Flags().Lookup(name))
}

func flagName(cmd *cobra.Command, name string) string {
	parentName := doctl.NSRoot
	if cmd.Parent() != nil {
		parentName = cmd.Parent().Name()
	}

	return fmt.Sprintf("%s-%s-%s", parentName, cmd.Name(), name)
}

func cmdNS(cmd *cobra.Command) string {
	parentName := doctl.NSRoot
	if cmd.Parent() != nil {
		parentName = cmd.Parent().Name()
	}

	return fmt.Sprintf("%s-%s", parentName, cmd.Name())
}

type cmdRunner func(ns string, out io.Writer) error

func cmdBuilder(cr cmdRunner, cliText, desc string, out io.Writer, aliases ...string) *cobra.Command {
	return &cobra.Command{
		Use:     cliText,
		Aliases: aliases,
		Short:   desc,
		Long:    desc,
		Run: func(cmd *cobra.Command, args []string) {
			checkErr(cr(cmdNS(cmd), out), cmd)
		},
	}
}

func listDroplets(client *godo.Client) ([]godo.Droplet, error) {
	f := func(opt *godo.ListOptions) ([]interface{}, *godo.Response, error) {
		list, resp, err := client.Droplets.List(opt)
		if err != nil {
			return nil, nil, err
		}

		si := make([]interface{}, len(list))
		for i := range list {
			si[i] = list[i]
		}

		return si, resp, err
	}

	si, err := doctl.PaginateResp(f)
	if err != nil {
		return nil, err
	}

	list := make([]godo.Droplet, len(si))
	for i := range si {
		list[i] = si[i].(godo.Droplet)
	}

	return list, nil
}

func extractDropletPublicIP(droplet *godo.Droplet) string {
	for _, in := range droplet.Networks.V4 {
		if in.Type == "public" {
			return in.IPAddress
		}
	}

	return ""

}
