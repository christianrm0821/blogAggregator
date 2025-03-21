package config

import (
	"encoding/json"
	"io"
	"os"
)

func Read() (*Config, error) {
	var ans Config
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	path := home + "/.gatorconfig.json"

	res, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	data, err := io.ReadAll(res)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &ans)
	if err != nil {
		return nil, err
	}

	return &ans, nil
}
