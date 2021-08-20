package cmd

import (
	"fmt"
	"github.com/StevenACoffman/jira-tool/pkg/atlassian"
	"os"

	"github.com/spf13/cobra"
)

// reviewCmd represents the review command
var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Transition an issue to Review status",
	Long: `Transition an issue to Review status`,
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

		err := atlassian.MoveIssueToStatus(jiraClient, issue, issueKey, atlassian.InReviewStatusID)
		if err != nil {
			fmt.Println(err)
			os.Exit(exitFail)
		}
		os.Exit(exitSuccess)
	},
}

func init() {
	rootCmd.AddCommand(reviewCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reviewCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reviewCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
