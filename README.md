# jt - jira-tool

`jt` is a CLI tool for viewing and manipulating JIRA issues.

### Usage:
jt [command]

### Available Commands:
| command | what it does |
|---|---|
| block       | Transition an issue to Blocked status |
| completion  | generate the autocompletion script for the specified shell |
| done        | Transition an issue to Deployed / Done status |
| help        | Help about any command |
| land        | Transition an issue to Landed status |
| onit        | Self-assign and transition an issue to In Progress status |
| review      | Transition an issue to Review status |
| take        | Assign an issue to you |
| todo        | Transition an issue to To Do status |
| triage      | Transition an issue to Triage status |
| wti         | What The Issue? - View an issue |

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

### Development and Limitations
Currently, the config does not allow overriding the workflow states.

Also, if a user doesn't have a config file, it should help them create one.

