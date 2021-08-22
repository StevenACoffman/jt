package cmd

import (
	"fmt"
	"github.com/StevenACoffman/jt/pkg/atlassian"
	"os"

	"github.com/spf13/cobra"
)

// todoCmd represents the todo command
var todoCmd = &cobra.Command{
	Use:   "todo",
	Short: "Transition an issue to To Do status",
	Long: `Transition an issue to To Do status`,
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

		err := atlassian.MoveIssueToStatus(jiraClient, issue, issueKey, atlassian.ToDoStatusID)
		if err != nil {
			fmt.Println(err)
			os.Exit(exitFail)
		}
		os.Exit(exitSuccess)
	},
}

func init() {
	rootCmd.AddCommand(todoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// todoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// todoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
