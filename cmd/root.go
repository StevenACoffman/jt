package cmd

import (
  "fmt"
  "github.com/StevenACoffman/jt/pkg/atlassian"
  "github.com/spf13/cobra"
  "os"

  "github.com/andygrunwald/go-jira"
  homedir "github.com/mitchellh/go-homedir"
  "github.com/spf13/viper"
)

const (
  // exitFail is the exit code if the program fails.
  exitFail = 1
  // exitSuccess is the exit code if the program succeeds
  exitSuccess = 0
)

var cfgFile string
var jiraClient *jira.Client


// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
  Use:   "jt",
  Short: "jt - JIRA Issue Tool",
  Long: `jt is a CLI tool for viewing and manipulating JIRA issues.`,
  Args: cobra.RangeArgs(1, 2),
  // Uncomment the following line if your bare application
  // has an action associated with it:
  	Run: func(cmd *cobra.Command, args []string) {
      if len(args) == 0 {
        fmt.Println("You need to pass a desired jira status argument (and maybe a jira issue like TEAM-1234)")
        os.Exit(exitFail)
      }
      var issueKey string
      statusName := args[0]
      if len(args) > 1 {
        issueKey = args[1]
      } else {

      }
      issue, _, issueErr := jiraClient.Issue.Get(issueKey, nil)
      if issueErr != nil {
        fmt.Printf("Unable to get Issue %s: %+v", issueKey, issueErr)
        os.Exit(exitFail)
      }

      err := atlassian.MoveIssueToStatusByName(jiraClient, issue, issueKey, statusName)
      if err != nil {
        fmt.Println(err)
        os.Exit(exitFail)
      }
      os.Exit(exitSuccess)
   },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func init() {
  cobra.OnInitialize(initConfig)

  // Here you will define your flags and configuration settings.
  // Cobra supports persistent flags, which, if defined here,
  // will be global for your application.

  rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/jira)")
}


// initConfig reads in config file and ENV variables if set.
func initConfig() {
  // default delimiter is "." and emails contain these
  v := viper.NewWithOptions(viper.KeyDelimiter("::"))
  v.SetConfigType("json")

  if cfgFile != "" {
    // Use config file from the flag.
    v.SetConfigFile(cfgFile)
  } else {
    // Find home directory.
    home, err := homedir.Dir()
    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }

    // Search config in home directory with name ".jt" (without extension).
    v.AddConfigPath(home+"/.config")
    v.SetConfigName("jira")
  }

  // If a config file is found, read it in.
  if err := v.ReadInConfig(); err != nil {
    fmt.Println("Unable to read config using config file:", v.ConfigFileUsed())
    return
  }

  jiraConfig := atlassian.Config{
    Token:  getEnv("ATLASSIAN_API_TOKEN", v.GetString("token")),
    User:  getEnv("ATLASSIAN_API_USER", v.GetString("user")),
    Host: getEnv("ATLASSIAN_HOST", v.GetString("host")),
  }
  jiraClient = atlassian.GetJIRAClient(&jiraConfig)

}

func getEnv(key, fallback string) string {
  if value, ok := os.LookupEnv(key); ok {
    return value
  }
  return fallback
}

