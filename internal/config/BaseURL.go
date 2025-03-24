package config

import "os"

func BaseURL() string {
	home, _ := os.UserHomeDir()
	baseURL := home + "/.gatorconfig.json"
	return baseURL

}
