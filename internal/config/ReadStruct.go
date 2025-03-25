package config

import (
	"encoding/json"
	"io"
	"os"
)

// Reads the data from the path "baseurl" and return a config pointer to that data or an error
func Read() (*Config, error) {
	var ans Config
	path := BaseURL()

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
