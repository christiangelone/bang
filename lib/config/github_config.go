package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type GitHub struct {
	GitHubAuthToken string `yaml:"GITHUB_AUTH_TOKEN"`
}

func GetGitHubConfig() GitHub {
	config := Config{}
	configPath, pathErr := getConfigPath()
	if pathErr != nil {
		return config.GitHub
	}
	if _, existErr := os.Stat(configPath); existErr == nil || !os.IsNotExist(existErr) {
		configFile, readErr := ioutil.ReadFile(configPath)
		if readErr != nil {
			return config.GitHub
		}

		yamlErr := yaml.Unmarshal(configFile, &config)
		if yamlErr != nil {
			return GitHub{}
		}
	}
	return config.GitHub
}
