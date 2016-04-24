package cmd

import "github.com/spf13/cobra"

func init() {
	OnsCmd.AddCommand(lsCmd)
}

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all DNS records of the zone",
	Run: func(cmd *cobra.Command, args []string) {

		records, err := onsClient.Ls(zone)
		if err != nil {
			exit("Fail to list records", err)
		}

		for _, record := range records {
			record.Print()
		}
	},
}
