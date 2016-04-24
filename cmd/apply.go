package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	OnsCmd.AddCommand(applyCmd)
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Changes DNS",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Refreshing DNS state prior to apply...\n\n")

		added, removed, err := onsClient.Apply(zone)
		if err != nil {
			exit("Fail to apply DNS configuration", err)
		}

		if (added + removed) > 0 {
			fmt.Println("")
		}
		cyan("Apply: %d added, %d removed.\n", added, removed)
	},
}
