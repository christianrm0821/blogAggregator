package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// converts data with new user to json and writes it onto the file in path "BaseURL()"
func (c *Config) SetUser(username *string) error {
	(*c).CurrentUserName = username

	data, err := json.MarshalIndent(*c, "", " ")
	if err != nil {
		return err
	}

	JsonFile := BaseURL()
	w, err := os.OpenFile(JsonFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("error making the writer")
		return err
	}
	defer w.Close()

	_, err = w.Write(data)
	if err != nil {
		fmt.Println("error in writing to file")
		return err
	}
	return nil

}
