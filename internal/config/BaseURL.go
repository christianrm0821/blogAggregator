package config

import "os"

// The base url to the json file on the computer
func BaseURL() string {
	home, _ := os.UserHomeDir()
	baseURL := home + "/.gatorconfig.json"
	return baseURL

}
