package git

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func command(out io.Writer, cmds []string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("git.exe", cmds...)
	case "linux", "darwin":
		cmd = exec.Command("git", cmds...)
	default:
		return fmt.Errorf("unsupported platform")
	}

	cmd.Stdin = os.Stdin
	if out != nil {
		cmd.Stdout = out
		cmd.Stderr = out
	}

	return cmd.Run()
}

func CurrentBranch() string {
	var buf bytes.Buffer
	err := command(&buf, []string{"symbolic-ref", "--short", "HEAD"})
	if err != nil {
		return ""
	}
	return strings.TrimSpace(buf.String())
}
