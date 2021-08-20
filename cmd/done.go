/*
Copyright © 2021 Steve Coffman <steve@khanacademy.org>

*/
package cmd

import (
	"fmt"
	"github.com/StevenACoffman/jira-tool/pkg/atlassian"
	"os"

	"github.com/spf13/cobra"
)

// doneCmd represents the done command
var doneCmd = &cobra.Command{
	Use:   "done",
	Short: "Transition an issue to Deployed / Done status",
	Long: `Transition an issue to Deployed / Done status`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("You failed to pass a jira issue argument")
			os.Exit(exitFail)
		}
		issueKey := args[0]
		issue, _, issueErr := jiraClient.Issue.Get(issueKey, nil)
		if issueErr != nil {
			fmt.Printf("Unable to get Issue %s: %+v", issueKey, issueErr)
			os.Exit(exitFail)
		}

		err := atlassian.MoveIssueToStatus(jiraClient, issue, issueKey, atlassian.DeployedDoneStatusID)
		if err != nil {
			fmt.Println(err)
			os.Exit(exitFail)
		}
		os.Exit(exitSuccess)
	},
}

func init() {
	rootCmd.AddCommand(doneCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// doneCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// doneCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
