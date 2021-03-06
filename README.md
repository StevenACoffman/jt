# jt - jira-tool

`jt` is a CLI tool for viewing and manipulating JIRA issues.

One common example usage to transition an issue to a new status:
```
jt "In Progress" TEAM-1234
```

If you are in a git repository where the topic branch's name matches `[whatever-]team-1234[-whatever]`, you can omit
the issue argument as it is implied.

Yeah, we even let you use underscores.

### Common Usage:
jt [new state] [issue number]

**Note:** 

We case insensitively look for valid transition states in your issue's workflow. If you give `tRiAgE`
we will find `Triage`, if that is a valid transition for your issue's current status.

If no valid transition state matches *exactly*, we then try matching against
possible states that have had their whitespace removed. If you give "todo" we will find possible state `To Do`.

If still no valid transition state is matched, we will then try partial match, so that
"done" will match possible state `Deployed / Done`.

This will otherwise only transition an issue to a matching valid state according to your
JIRA board's workflow.

### Other Available Commands:
| command | what it does |
|---|---|
| onit        | Self-assign and transition an issue to In Progress status |
| take        | Assign an issue to you |
| wti         | What The Issue? - View an issue in Github Markdown |
| config      | Will save the JIRA token, email, and tenant url to a config file
| completion  | generate the autocompletion script for the specified shell |
| help        | Help about any command |

Shared Flags:
| flag | what it does |
|---|---|
| --config string |  config file (default is $HOME/.config/jira) |
| -h, --help      |  help for jt |

### Tips
Use "jt [command] --help" for more information about a command.

### Installation
Homebrew users can do this:
```
brew tap StevenACoffman/jt
brew install jt
```

Go developers with `$HOME/bin` in their `$PATH` can run `mage` if they have [mage](https://magefile.org/) installed.

Alternatively, `go run mage.go` will work even without `mage` installed, but it will still put the binary in `$HOME/bin`. 

### Future Development
Currently, this tool is all JIRA, other than translating Atlassian Markup into Github Markdown.

I would like to combine JIRA + Github steps. For instance, moving a story to `In Progress`,
while also creating a Github Pull Request in Draft. 

It would be nice if it took the JIRA story's description's translation as Github Markdown
and used it in creating the Github Pull Request. 

The programming is easy, but it is tricky to find unambiguous commands that fit all possible workflows.

### Alternatives

There is another [jira cli](https://github.com/go-jira/jira) that is quite sophisticated, featureful,
and maybe complicated, but I found custom workflow transitions either didn't work, or were cumbersome.

### Credit Where Credit is Due

This tool uses some excellent libraries and tools like [andygrunwald/go-jira](https://github.com/andygrunwald/go-jira),
[spf13/cobra](https://github.com/spf13/cobra), [mage](https://magefile.org/), [moul/http2curl](https://github.com/moul/http2curl) and [bubbletea](github.com/charmbracelet/bubbletea)