package cmd

import (
	"fmt"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thbkrkr/ons/client"
)

// OnsCmd is the ONS CLI main command
var OnsCmd = &cobra.Command{
	Use:   "ons",
	Short: "Utility to manage OVH DNS zone.",
}

const (
	envPrefix = "ons"
)

var (
	onsClient *client.OnsClient
	zone      string

	onsDir     string
	statePath  string
	configPath string

	magenta = color.New(color.FgMagenta).SprintFunc()
	green   = color.New(color.FgGreen).SprintFunc()
	cyan    = color.New(color.FgCyan).PrintfFunc()
)

func init() {
	viper.SetEnvPrefix("ons")
	viper.AutomaticEnv()

	viper.SetDefault("path", "dns")
	viper.SetDefault("endpoint", "ovh-eu")

	zone = env("zone")
	onsDir = env("path")

	statePath = onsDir + "/ons.state.json"
	configPath = onsDir + "/ons.config.json"

	var err error
	onsClient, err = client.NewOnsClient(statePath, configPath,
		env("endpoint"), env("ak"), env("as"), env("ck"))

	if err != nil {
		exit("Fail to start ons", err)
	}
}

func env(key string) string {
	value := viper.GetString(key)
	if value == "" {
		log.Errorf("%s not defined.", strings.ToUpper(envPrefix+"_"+key))
		os.Exit(1)
	}
	return value
}

func require(cmd string, min int, max int, args []string) {
	if len(args) < min || len(args) > max {
		exit(fmt.Sprintf("`%s` requires %d argument.", cmd, max), nil)
	}
}

func exit(msg string, err error) {
	fmt.Println()
	if err != nil {
		log.WithError(err).Error(msg)
	} else {
		log.Error(msg)
	}

	os.Exit(1)
}
