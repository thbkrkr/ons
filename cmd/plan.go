package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	printAddition = color.New(color.FgGreen).PrintfFunc()
	printRemoval  = color.New(color.FgRed).PrintfFunc()
)

func init() {
	OnsCmd.AddCommand(planCmd)
}

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Show the execution plan",
	Run: func(cmd *cobra.Command, args []string) {

		plan()
	},
}

func plan() {
	fmt.Printf("Refreshing DNS zone state prior to plan...\n\n")

	toAdd, toRm, err := onsClient.Plan(zone)
	if err != nil {
		exit("Fail to plan", err)
	}

	for _, r := range toAdd {
		printAddition("+ dns record: %-16s %s.%s\n", r.Target, r.SubDomain, zone)
	}
	for _, r := range toRm {
		comment := ""
		if r.ID == 0 {
			comment = "(already removed from the DNS zone)"
		}
		printRemoval("- dns record: %-16s %s.%s %s\n", r.Target, r.SubDomain, zone, comment)
	}

	if len(toAdd)+len(toRm) > 0 {
		fmt.Println()
	}

	cyan("Plan: %d to add, %d to remove.\n", len(toAdd), len(toRm))
}
