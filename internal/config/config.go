package config 

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	configFileName = ".gatorconfig.json"
)

type Config struct{
	Db_url string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func ReadJson() (Config, error){
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	content, err := os.ReadFile(configFilePath)
	if err != nil {
		return Config{}, err
	}
	
	var config Config
	if err = json.Unmarshal(content, &config); err != nil {
		return Config{}, err
	}
	return config, nil
}

func (c *Config) SetUser(user_name string) error {
	c.Current_user_name = user_name
	return writeJson(c,)
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, configFileName), nil 
}

func writeJson(cfg *Config) error {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	
	jsonBytes, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	err = os.WriteFile(configFilePath, jsonBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}