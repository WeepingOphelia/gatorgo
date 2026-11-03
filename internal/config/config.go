package config

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	Url               string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

const cfgFilename = ".gatorconfig.json"

func Read() (Config, error) {
	path, err := getCfgFilePath()
	if err != nil {
		return Config{}, err
	}
	f, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	dat, err := io.ReadAll(f)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	err = json.Unmarshal(dat, &cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (c Config) SetUser(user_name string) error {
	c.Current_user_name = user_name
	return write(c)
}

func getCfgFilePath() (string, error) {
	path, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return path + "/" + cfgFilename, nil
}

func write(cfg Config) error {
	path, err := getCfgFilePath()
	if err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	dat, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	_, err = f.Write(dat)
	if err != nil {
		return err
	}
	return nil
}
