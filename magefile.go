// +build mage

package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh"
)

var (
	// Default target to run when none is specified
	// If not set, running mage will list available targets
	Default        = Install
	appName = "jt"
	installPath    string
	currentWorkDir string
	// allow user to override go executable by running as GOEXE=xxx make ... on unix-like systems
	goexe = "go"
)

// A build step that requires additional params, or platform specific steps for example
func Build() error {
	fmt.Println("Building...")
	cmd := exec.Command("go", "build", "-o", appName, currentWorkDir)
	return cmd.Run()
}

// A custom install step if you need your bin someplace other than go/bin
func Install() error {
	mg.Deps(Build)
	fmt.Println("Installing...")
	return os.Rename(filepath.Join(currentWorkDir,"jt"), installPath)
}

// Clean up after yourself
func Clean() {
	fmt.Println("Cleaning...")
	os.RemoveAll("jt")
}

var releaseTag = regexp.MustCompile(`^v[0-9]+\.[0-9]+\.[0-9]+$`)

// Generates a new release. Expects a version tag in vx.x.x format.
// really only useful for maintainers
func Release(tag string) (err error) {
	if !releaseTag.MatchString(tag) {
		return errors.New("TAG environment variable must be in semver vx.x.x format, but was " + tag)
	}

	if err := sh.RunV("git", "tag", "-a", tag, "-m", tag); err != nil {
		return err
	}
	if err := sh.RunV("git", "push", "origin", tag); err != nil {
		return err
	}
	defer func() {
		if err != nil {
			sh.RunV("git", "tag", "--delete", "$TAG")
			sh.RunV("git", "push", "--delete", "origin", "$TAG")
		}
	}()
	return sh.RunV("goreleaser", "--rm-dist")
}


// tag returns the git tag for the current branch or "" if none.
func tag() string {
	s, _ := sh.Output("git", "describe", "--tags")
	return s
}

// hash returns the git hash for the current repo or "" if none.
func hash() string {
	hash, _ := sh.Output("git", "rev-parse", "--short", "HEAD")
	return hash
}


func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// this is probably overkill, but in case they are not Go developers
func init() {
	if exe := os.Getenv("GOEXE"); exe != "" {
		goexe = exe
	}

	// We want to use Go 1.11 modules even if the source lives inside GOPATH.
	// The default is "auto".
	os.Setenv("GO111MODULE", "on")

	os.Setenv("APP", appName)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	currentWorkDir, err = os.Getwd()
	if err != nil {
		panic(err)
	}
	goBinDir := getEnv("GOBIN", fmt.Sprintf("%s/go/bin", homeDir))
	os.Setenv("GOBIN", goBinDir)
	// mkdir -p ${HOME}/go/bin
	err = os.MkdirAll(goBinDir, 0o755)
	if err != nil {
		panic(err)
	}
	installPath = fmt.Sprintf("%s/%s", goBinDir, appName)
	os.Setenv("INSTALLPATH", installPath)
	os.Setenv("PATH", fmt.Sprintf("%s:%s", goBinDir, os.Getenv("PATH")))
	os.Setenv("GOPRIVATE", "github.com/Khan")
}