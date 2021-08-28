package atlassian

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/StevenACoffman/jt/pkg/middleware"

	"github.com/andygrunwald/go-jira"
)

// GetJIRAClient takes a config, and makes a JIRAClient configured
// to use BasicAuth
func GetJIRAClient(config *Config) *jira.Client {
	httpClient := middleware.NewBasicAuthHTTPClient(config.User, config.Token)

	jiraClient, err := jira.NewClient(httpClient, config.Host)
	if err != nil {
		log.Fatalf("unable to create new JIRA client. %v", err)
	}
	return jiraClient
}

// GetIssue checks if issue exists in the JIRA instance.
// If not an error will be returned.
func GetIssue(jiraClient *jira.Client, issue string) (*jira.Issue, error) {
	jiraIssue, resp, err := jiraClient.Issue.Get(issue, nil)
	if err != nil {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("got empty response from jira for issue %s", issue)
	}
	if c := resp.StatusCode; c < 200 || c > 299 {
		return nil, fmt.Errorf(
			"jira Request for issue %s returned %s (%d)",
			issue,
			resp.Status,
			resp.StatusCode,
		)
	}
	return jiraIssue, nil
}

// Jiration - convenience for Jira Markup to Github Markdown translation rule
type Jiration struct {
	re   *regexp.Regexp
	repl interface{}
}

// JiraToMD - This uses some regular expressions to make a reasonable translation
// from Jira Markup to Github Markdown. It is not a complete PEG so it will break down
// especially for more complicated nested formatting (lists inside of lists)
func JiraToMD(str string) string {
	jirations := []Jiration{
		{ // UnOrdered Lists
			re: regexp.MustCompile(`(?m)^[ \t]*(\*+)\s+`),
			repl: func(groups []string) string {
				_, stars := groups[0], groups[1]
				return strings.Repeat("  ", len(stars)-1) + "* "
			},
		},
		{ // Ordered Lists
			re: regexp.MustCompile(`(?m)^[ \t]*(#+)\s+`),
			repl: func(groups []string) string {
				_, nums := groups[0], groups[1]
				return strings.Repeat("  ", len(nums)-1) + "1. "
			},
		},
		{ // Headers 1-6
			re: regexp.MustCompile(`(?m)^h([0-6])\.(.*)$`),
			repl: func(groups []string) string {
				_, level, content := groups[0], groups[1], groups[2]
				i, _ := strconv.Atoi(level)
				return strings.Repeat("#", i) + content
			},
		},
		{ // Bold
			re:   regexp.MustCompile(`\*(\S.*)\*`),
			repl: "**$1**",
		},
		{ // Italic
			re:   regexp.MustCompile(`\_(\S.*)\_`),
			repl: "*$1*",
		},
		{ // Monospaced text
			re:   regexp.MustCompile(`\{\{([^}]+)\}\}`),
			repl: "`$1`",
		},
		{ // Citations (buggy)
			re:   regexp.MustCompile(`\?\?((?:.[^?]|[^?].)+)\?\?`),
			repl: "<cite>$1</cite>",
		},
		{ // Inserts
			re:   regexp.MustCompile(`\+([^+]*)\+`),
			repl: "<ins>$1</ins>",
		},
		{ // Superscript
			re:   regexp.MustCompile(`\^([^^]*)\^`),
			repl: "<sup>$1</sup>",
		},
		{ // Subscript
			re:   regexp.MustCompile(`~([^~]*)~`),
			repl: "<sub>$1</sub>",
		},
		{ // Strikethrough
			re:   regexp.MustCompile(`(\s+)-(\S+.*?\S)-(\s+)`),
			repl: "$1~~$2~~$3",
		},
		{ // Code Block
			re: regexp.MustCompile(
				`\{code(:([a-z]+))?([:|]?(title|borderStyle|borderColor|borderWidth|bgColor|titleBGColor)=.+?)*\}`,
			),
			repl: "```$2",
		},
		{ // Code Block End
			re:   regexp.MustCompile(`{code}`),
			repl: "```",
		},
		{ // Pre-formatted text
			re:   regexp.MustCompile(`{noformat}`),
			repl: "```",
		},
		{ // Un-named Links
			re:   regexp.MustCompile(`(?U)\[([^|]+)\]`),
			repl: "<$1>",
		},
		{ // Images
			re:   regexp.MustCompile(`!(.+)!`),
			repl: "![]($1)",
		},
		{ // Named Links
			re:   regexp.MustCompile(`\[(.+?)\|(.+)\]`),
			repl: "[$1]($2)",
		},
		{ // Single Paragraph Blockquote
			re:   regexp.MustCompile(`(?m)^bq\.\s+`),
			repl: "> ",
		},
		{ // Remove color: unsupported in md
			re:   regexp.MustCompile(`(?m)\{color:[^}]+\}(.*)\{color\}`),
			repl: "$1",
		},
		{ // panel into table
			re: regexp.MustCompile(
				`(?m)\{panel:title=([^}]*)\}\n?(.*?)\n?\{panel\}`,
			),
			repl: "\n| $1 |\n| --- |\n| $2 |",
		},
		{ // table header
			re: regexp.MustCompile(`(?m)^[ \t]*((?:\|\|.*?)+\|\|)[ \t]*$`),
			repl: func(groups []string) string {
				_, headers := groups[0], groups[1]
				reBarred := regexp.MustCompile(`\|\|`)

				singleBarred := reBarred.ReplaceAllString(headers, "|")
				fillerRe := regexp.MustCompile(`\|[^|]+`)
				return "\n" + singleBarred + "\n" + fillerRe.ReplaceAllString(
					singleBarred,
					"| --- ",
				)
			},
		},
		{ // remove leading-space of table headers and rows
			re:   regexp.MustCompile(`(?m)^[ \t]*\|`),
			repl: "|",
		},
	}
	for _, jiration := range jirations {
		switch v := jiration.repl.(type) {
		case string:
			str = jiration.re.ReplaceAllString(str, v)
		case func([]string) string:
			str = ReplaceAllStringSubmatchFunc(jiration.re, str, v)
		default:
			fmt.Printf("I don't know about type %T!\n", v)
		}
	}
	return str
}

type JiraResolver struct {
	JiraClient *jira.Client
}

// JiraMarkupMentionToEmail will replace JiraMarkup account mentions
// with Display Name followed by parenthetical email addresses
func (j *JiraResolver) JiraMarkupMentionToEmail(str string) string {
	re := regexp.MustCompile(`(?m)(\[~accountid:)([a-zA-Z0-9-:]+)(\])`)
	rfunc := func(groups []string) string {
		// groups[0] is initial match
		accountID := groups[2]

		jiraUser, resp, err := j.JiraClient.User.Get(accountID)
		// if we cannot resolve it, so just leave it as it was
		if err != nil {
			return groups[0]
		}
		if resp == nil {
			return groups[0]
		}
		if c := resp.StatusCode; c < 200 || c > 299 {
			return groups[0]
		}

		return DisplayJiraUser(jiraUser)
	}
	return ReplaceAllStringSubmatchFunc(re, str, rfunc)
}

func DisplayJiraUser(jiraUser *jira.User) string {
	return jiraUser.DisplayName + " (" + jiraUser.EmailAddress + ")"
}

// ReplaceAllStringSubmatchFunc - Invokes Callback for Regex Replacement
// The repl function takes an unusual string slice argument:
// - The 0th element is the complete match
// - The following slice elements are the nth string found
// by a parenthesized capture group (including named capturing groups)
//
// This is a Go implementation to match other languages:
// PHP: preg_replace_callback($pattern, $callback, $subject)
// Ruby: subject.gsub(pattern) {|match| callback}
// Python: re.sub(pattern, callback, subject)
// JavaScript: subject.replace(pattern, callback)
// See https://gist.github.com/elliotchance/d419395aa776d632d897
func ReplaceAllStringSubmatchFunc(
	re *regexp.Regexp,
	str string,
	repl func([]string) string,
) string {
	result := ""
	lastIndex := 0

	for _, v := range re.FindAllSubmatchIndex([]byte(str), -1) {
		groups := []string{}
		for i := 0; i < len(v); i += 2 {
			if v[i] == -1 || v[i+1] == -1 {
				// if the group is not found, avoid possible error
				groups = append(groups, "")
			} else {
				groups = append(groups, str[v[i]:v[i+1]])
			}
		}

		result += str[lastIndex:v[0]] + repl(groups)
		lastIndex = v[1]
	}

	return result + str[lastIndex:]
}

func JiraMarkupToGithubMarkdown(jiraClient *jira.Client, str string) string {
	jiraAccountResolver := JiraResolver{
		JiraClient: jiraClient,
	}
	resolved := jiraAccountResolver.JiraMarkupMentionToEmail(str)
	return JiraToMD(resolved)
}

func AssignIssueToSelf(jiraClient *jira.Client, issue *jira.Issue, issueKey string) error {
	self, _, selfErr := jiraClient.User.GetSelf()
	if selfErr != nil {
		return fmt.Errorf("unable to get myself: %+v", selfErr)
	}

	if issue.Fields.Assignee == nil || self.AccountID != issue.Fields.Assignee.AccountID {
		_, assignErr := jiraClient.Issue.UpdateAssignee(issueKey, self)
		if assignErr != nil {
			return fmt.Errorf("unable to assign %s to yourself: %+v", issueKey, assignErr)
		}
		fmt.Printf("Re-Assigned %s from %s\n", issueKey, DisplayJiraUser(issue.Fields.Assignee))
	} else {
		fmt.Println("Already assigned to to you")
	}
	return nil
}

// ParseJiraIssue - Sanitizes input
//  + Trims leading and trailing whitespace
//  + Trims a browse URL
//  + Trims anything after ABCD-1234
// If there is no jira issue match, returns empty string
func ParseJiraIssue(issueKey, host string) string {
	issueKey = strings.TrimSpace(issueKey)
	if issueKey == "" {
		return issueKey
	}
	if strings.HasPrefix(issueKey, host) {
		issueKey = strings.TrimPrefix(issueKey, host)
		issueKey = strings.TrimPrefix(issueKey, "/browse/")
	}
	// This will remove everything after the ABCD-1234
	reg := regexp.MustCompile(`(.*/)?(?P<Jira>[A-Za-z]+-[0-9]+).*`)
	if reg.MatchString(issueKey) {
		res := reg.ReplaceAllString(issueKey, "${Jira}")
		return res
	}
	return ""
}

// ParseJiraIssueFromBranch - Sanitizes input
//  + Trims leading "feature/" (or whatever GIT_WORKON_PREFIX set to)
//  + Trims leading and trailing whitespace
//  + Trims a browse URL
//  + Trims anything after ABCD-1234
// If there is no jira issue match, returns whatever was passed
func ParseJiraIssueFromBranch(issueKey, host, branchPrefix string) string {
	issueKey = strings.TrimSpace(issueKey)
	issueKey = strings.TrimPrefix(issueKey, branchPrefix)
	if issueKey == "" {
		return issueKey
	}
	if strings.HasPrefix(issueKey, host) {
		issueKey = strings.TrimPrefix(issueKey, host)
		issueKey = strings.TrimPrefix(issueKey, "/browse/")
	}

	res := TrimJira(issueKey)
	if res != "" {
		res = fmt.Sprintf("https://khanacademy.atlassian.net/browse/%s", res)
	}
	return res
}

// TrimJira will remove everything before and after the last ABCD-1234
// returns empty string if no jira issue is found
func TrimJira(s string) string {
	var re *regexp.Regexp
	var result string
	re = regexp.MustCompile("([a-zA-Z]{1,4}-[1-9][0-9]{0,6})")
	matches := re.FindAllStringSubmatch(s, -1)
	for _, match := range matches {
		for _, m := range match {
			result = m
		}
	}
	return result
}

func MoveIssueToStatusByName(jiraClient *jira.Client, issue *jira.Issue, issueKey string, statusName string) error {
	originalStatus := issue.Fields.Status.Name
	if issue.Fields.Status.Name == statusName ||
		CaseInsensitiveContains(issue.Fields.Status.Name, statusName){
		return fmt.Errorf("issue is Already in Status %s\n", issue.Fields.Status.Name)
	}

	err := transitionIssueByStatusName(jiraClient, issueKey, statusName)
	if err != nil {
		return err
	}
	issue, _, err = jiraClient.Issue.Get(issueKey, nil)
	if err != nil {
		return err
	}
	fmt.Printf("Issue %s Status successfully changed from: %s and set to: %+v\n",
		issueKey, originalStatus, issue.Fields.Status.Name)

	return nil
}

func transitionIssueByStatusName(jiraClient *jira.Client, issueKey string, statusName string) error {
	var transitionID string
	possibleTransitions, _, err := jiraClient.Issue.GetTransitions(issueKey)
	if err != nil {
		return err
	}
	for _, v := range possibleTransitions {
		if strings.EqualFold(v.To.Name, statusName) {
			transitionID = v.ID
			break
		}
	}
	// no exact match, so remove whitespace so that "ToDo" arg will match "TO DO" status
	if transitionID == "" {
		for _, v := range possibleTransitions {
			if strings.EqualFold(RemoveWhiteSpace(v.To.Name), statusName) {
				transitionID = v.ID
				break
			}
		}
	}
	// still no match, so look for partial, so "Done" arg will match "Deployed / Done"
	if transitionID == "" {
		// substring match only if exact match fails
		for _, v := range possibleTransitions {
			if CaseInsensitiveContains(v.To.Name, statusName) {
				transitionID = v.ID
				break
			}
		}
	}

	if transitionID == "" {
		return fmt.Errorf("there does not appear to be a valid transition to %s", statusName)
	}
	_, err = jiraClient.Issue.DoTransition(issueKey, transitionID)
	return err
}

func RemoveWhiteSpace(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for _, ch := range str {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
		}
	}
	return b.String()
}

func CaseInsensitiveContains(s, substr string) bool {
	s, substr = strings.ToUpper(s), strings.ToUpper(substr)
	return strings.Contains(s, substr)
}