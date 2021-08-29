package cmd

import (
	"fmt"
	"os"

	"github.com/StevenACoffman/jt/pkg/atlassian"

	"github.com/spf13/cobra"
)

// takeCmd represents the take command
var takeCmd = &cobra.Command{
	Use:   "take",
	Short: "Assign an issue to you",
	Long:  `Assign an issue to you`,
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		if jiraConfig == nil {
			configure()
		}
		var issueKey string
		if len(args) == 0 {
			issueKey = getIssueFromGitBranch()
		} else {
			issueKey = args[0]
		}
		issue, _, issueErr := jiraClient.Issue.Get(issueKey, nil)
		if issueErr != nil {
			fmt.Printf("Unable to get Issue %s: %+v", issueKey, issueErr)
			os.Exit(exitFail)
		}
		err := atlassian.AssignIssueToSelf(jiraClient, issue, issueKey)
		if err != nil {
			fmt.Println(err)
			os.Exit(exitFail)
		}
	},
}

func init() {
	rootCmd.AddCommand(takeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// takeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// takeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
