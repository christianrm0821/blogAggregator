package config

import "os"

// The base url to the json file on the computer
func BaseURL() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	baseURL := home + "/.gatorconfig.json"
	return baseURL, nil

}
