package config

import (
	"gopkg.in/yaml.v2"
	"io"
	"os"
)

// Define the structure for the YAML configuration
type Config struct {
	Targets []Target `yaml:"targets"`
	Server  Server   `yaml:"server"`
}

type Target struct {
	Name     string `yaml:"name"`
	HostAddr string `yaml:"host_addr"`
	Path     string `yaml:"path"`
}

type Server struct {
	Port                  int    `yaml:"port"`
	Path                  string `yaml:"path"`
	Verbose               bool   `yaml:"verbose"`
	TlsInsecureSkipVerify bool   `yaml:"tls_insecure_skip_verify"`
}

// GetConfig reads the YAML file and parses it into a Config struct
func GetConfig(filename string) (Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	return parseConfig(file)
}

// parseConfig reads the content of the file and unmarshals it into a Config struct
func parseConfig(reader io.Reader) (Config, error) {
	yamlFile, err := io.ReadAll(reader)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
