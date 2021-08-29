package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/StevenACoffman/jt/pkg/atlassian"
	"github.com/StevenACoffman/jt/pkg/colors"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Save your JIRA config for use in other commands",
	Long: `This will ask for your JIRA token, tenant URL and email.
It will backup any existing config file and make a new one.`,
	Run: func(cmd *cobra.Command, args []string) {
		configure()
		os.Exit(exitSuccess)
	},
}

func configure() {
	if atlassian.CheckConfigFileExists(cfgFile) {

		backupErr := BackupConfigFile(cfgFile)
		if backupErr != nil {
			fmt.Println("Unable to backup config file!")
			os.Exit(exitFail)
		}
	}
	model := initialModel()

	if err := tea.NewProgram(&model).Start(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}

	err := atlassian.SaveConfig(cfgFile, jiraConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(exitFail)
	}
	jiraClient = atlassian.GetJIRAClient(jiraConfig)
	fmt.Println("Successfully wrote config to ", cfgFile)
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

var (
	focusedStyle = lipgloss.NewStyle().Foreground(
		lipgloss.AdaptiveColor{
			Light: colors.ANSIGreen.String(),
			Dark:  colors.ANSIBrightGreen.String(),
		})
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240")) //#585858
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244")) //#808080

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type model struct {
	focusIndex int
	inputs     []textinput.Model
	cursorMode textinput.CursorMode
	choice     chan *atlassian.Config
}

func initialModel() model {
	m := model{
		inputs: make([]textinput.Model, 3),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.NewModel()
		t.CursorStyle = cursorStyle
		t.CharLimit = 128

		switch i {
		case 0:
			t.Placeholder = "Paste token here"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		case 1:
			t.Placeholder = "Host URL like https://tenant.atlassian.net"
		case 2:
			t.Placeholder = "Email"
		}

		m.inputs[i] = t
	}

	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		// Change cursor mode
		case "ctrl+r":
			m.cursorMode++
			if m.cursorMode > textinput.CursorHide {
				m.cursorMode = textinput.CursorBlink
			}
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				cmds[i] = m.inputs[i].SetCursorMode(m.cursorMode)
			}
			return m, tea.Batch(cmds...)

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, save choices and exit.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				if jiraConfig == nil {
					jiraConfig = &atlassian.Config{}
				}
				for i, input := range m.inputs {
					switch i {
					case 0:
						jiraConfig.Token = input.Value()
					case 1:
						jiraConfig.Host = input.Value()
					case 2:
						jiraConfig.User = input.Value()
					}
				}
				return m, tea.Quit
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m model) View() string {
	var b strings.Builder
	b.WriteString("It looks like we need a Jira API Token.\n\n")
	styledLink := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.AdaptiveColor{Light: "4", Dark: "12"}). // Dark Blue or LightBlue
		Underline(true).
		Render("https://id.atlassian.com/manage/api-tokens")
	b.WriteString(fmt.Sprintf(
		"First, go to %s to create a personal api token.\n",
		styledLink))

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(helpStyle.Render("cursor mode is "))
	b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

	return b.String()
}

func BackupConfigFile(filename string) error {
	return os.Rename(filename, filename+".bak")
}
