package config

import (
	"os"
	"path/filepath"

	"github.com/christiangelone/bang/lib/file"
	"github.com/christiangelone/bang/lib/system"
	"gopkg.in/yaml.v2"
)

type Config struct {
	GitHub GitHub `yaml:"github"`
}

func New() Config {
	return Config{
		GitHub: GitHub{
			GitHubAuthToken: os.Getenv("GITHUB_AUTH_TOKEN"),
		},
	}
}

func getConfigPath() (string, error) {
	bangPath, err := system.BangFolderPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(bangPath, "config.yml"), nil
}

func NewIfNotExist() error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fileHandler := file.NewHandler()

		file, fileErr := fileHandler.OpenFile(configPath, os.O_CREATE|os.O_WRONLY, 0644)
		if fileErr != nil {
			return fileErr
		}

		config := New()
		configBytes, yamlErr := yaml.Marshal(&config)
		if yamlErr != nil {
			return yamlErr
		}

		_, writeErr := file.Write(configBytes)
		if writeErr != nil {
			return writeErr
		}
	}
	return nil
}
