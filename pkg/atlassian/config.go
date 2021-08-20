package atlassian

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// Config struct
type Config struct {
	Host  string `json:"host",mapstructure:"host"`
	User  string `json:"user",mapstructure:"user"`
	Token string `json:"token",mapstructure:"token"`
}

// ReadConfigFromFile returns an error if file does not exist
func ReadConfigFromFile() (*Config, error) {
	configFile, configErr := expandTilde(getEnv("ATLASSIAN_CONFIG_FILE", "~/.config/jira"))

	if configErr != nil {
		// if we can't get the config file, then we have no hope.
		return nil, fmt.Errorf("unable to get config file directory %+v", configErr)
	}

	var config Config
	configJSON, err := ioutil.ReadFile(configFile)
	if err != nil {
		return &config, err
	}

	err = json.Unmarshal(configJSON, &config)

	if err != nil {
		return &config, err
	}

	config.Token = getEnv("ATLASSIAN_API_TOKEN", config.Token)
	config.Host = getEnv("ATLASSIAN_HOST", config.Host)
	config.User = getEnv("ATLASSIAN_API_USER", config.User)

	return &config, nil
}

func ReadConfigFromEnv() *Config {
	host := os.Getenv("ATLASSIAN_HOST")
	username := os.Getenv("ATLASSIAN_API_USER")
	token := os.Getenv("ATLASSIAN_API_TOKEN")
	config := Config{
		Host:  host,
		User:  username,
		Token: token,
	}
	return &config
}

func ConfigureJira() *Config {
	config, err := ReadConfigFromFile()
	if err != nil {
		// we got an error reading from the config file, so just use env
		return ReadConfigFromEnv()
	}
	// allow env to replace file config
	config.Token = getEnv("ATLASSIAN_API_TOKEN", config.Token)
	config.Host = getEnv("ATLASSIAN_HOST", config.Host)
	config.User = getEnv("ATLASSIAN_API_USER", config.User)
	return config
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// "~/.gitignore" -> "/home/tyru/.gitignore"
func expandTilde(path string) (string, error) {
	if !strings.HasPrefix(path, "~") {
		return path, nil
	}
	var paths []string
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	for _, p := range strings.Split(path, string(filepath.Separator)) {
		if p == "~" {
			paths = append(paths, u.HomeDir)
		} else {
			paths = append(paths, p)
		}
	}
	return "/" + filepath.Join(paths...), nil
}
