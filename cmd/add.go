package cmd

import (
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	OnsCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add [subdomain] [ip]",
	Short: "Plan to add a record",
	Long:  "Plan to add a DNS zone record given a sub domain and an IP. If the IP is not set DOCKER_MACHINE_NAME is used and the IP is resolved using docker machine",
	Run: func(cmd *cobra.Command, args []string) {

		require("add", 1, 2, args)
		subDomain := args[0]
		target := argTarget(args)

		err := onsClient.Add(zone, subDomain, target)
		if err != nil {
			exit("Fail to add record", err)
		}

		plan()
	},
}

func argTarget(args []string) string {
	target := ""

	machine := os.Getenv("DOCKER_MACHINE_NAME")
	if machine != "" {
		output, err := exec.Command("docker-machine", "ip", machine).Output()
		if err != nil {
			exit("Fail to get ip from `"+machine+"` using docker-machine", nil)
		}
		target = strings.Replace(string(output), "\n", "", -1)
	} else {
		if len(args) == 2 {
			target = args[1]
		}
	}

	if target == "" {
		exit("`add` requires an ip argument or the environment variable DOCKER_MACHINE_NAME", nil)
	}

	return target
}
