package cmd

import (
	"fmt"
	"github.com/StevenACoffman/jt/pkg/atlassian"
	"os"

	"github.com/spf13/cobra"
)

// onitCmd represents the onit command
var onitCmd = &cobra.Command{
	Use:   "onit",
	Short: "Self-assign and transition an issue to In Progress status",
	Long: `Assign the issue to yourself and transition an issue to In Progress status`,
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

		err := atlassian.MoveIssueToStatus(jiraClient, issue, issueKey, atlassian.InProgressStatusID)
		if err != nil {
			fmt.Println(err)
			os.Exit(exitFail)
		}

		err = atlassian.AssignIssueToSelf(jiraClient, issue, issueKey)
		if err != nil {
			fmt.Println(err)
			os.Exit(exitFail)
		}

		os.Exit(exitSuccess)
	},
}

func init() {
	rootCmd.AddCommand(onitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// onitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// onitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
