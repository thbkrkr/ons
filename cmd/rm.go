package cmd

import "github.com/spf13/cobra"

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Plan to remove records matching a sub domain",
	Run: func(cmd *cobra.Command, args []string) {

		require("rm", 1, 1, args)
		subDomain := args[0]

		err := onsClient.Rm(zone, subDomain)
		if err != nil {
			exit("Fail to remove record", err)
		}

		plan()
	},
}

func init() {
	OnsCmd.AddCommand(rmCmd)
}
