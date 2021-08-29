package cmd

import (
	"fmt"

	"github.com/StevenACoffman/jt/pkg/atlassian"

	"github.com/spf13/cobra"
)

var omitTitle, omitDescription bool

// wtiCmd represents the wti command
var wtiCmd = &cobra.Command{
	Use:   "wti",
	Short: "What The Issue? - View an issue",
	Long:  `What The Issue? Will View an issue.`,
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

		jiraIssue, issueErr := atlassian.GetIssue(jiraClient, issueKey)
		if issueErr != nil {
			fmt.Println(issueErr)
		}

		if issueErr == nil && jiraIssue != nil {
			if !omitTitle {
				fmt.Printf("%s - %s\n\n", jiraIssue.Key, jiraIssue.Fields.Summary)
			}
			if !omitDescription {
				fmt.Println(
					atlassian.JiraMarkupToGithubMarkdown(
						jiraClient, jiraIssue.Fields.Description))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(wtiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// wtiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	flags := wtiCmd.Flags()
	//.BoolP("toggle", "t", false, "Help message for toggle")
	flags.BoolVarP(&omitTitle, "no-title", "t", false, "Do Not Print Title")
	flags.BoolVarP(&omitDescription, "no-description", "d", false, "Do Not Print Description")
}
