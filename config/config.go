package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

type configStruct struct {
	Host  string `json:"host"`
	Token string `json:"api-key"`
}

func parseConfig(filePath string) (*configStruct, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	cfg := &configStruct{}
	err = json.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func getExecDir() string {
	ex, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Dir(ex)
}

func Read() (*configStruct, error) {
	var configPath string

	if len(os.Args) >= 3 && os.Args[1] == "-c" {
		configPath = os.Args[2]
	}

	if configPath != "" {
		cfg, err := parseConfig(configPath)
		if err == nil {
			return cfg, err
		}
	}

	dir, _ := os.Getwd()

	configPath = path.Join(dir, "config.json")
	cfg, err := parseConfig(configPath)
	if err == nil {
		return cfg, err
	}

	configPath = path.Join(getExecDir(), "config.json")
	return parseConfig(configPath)
}
