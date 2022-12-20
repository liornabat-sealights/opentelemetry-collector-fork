package slauth

import (
	"encoding/json"
	"io/ioutil"
)

type Settings struct {
	AuthEndpoint     string `json:"authEndpoint"`
	AuthVerbEndpoint string `json:"authVerbEndpoint"`
}

func NewSettings() *Settings {
	return &Settings{}
}

func (s *Settings) LoadSettings(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, s)
	if err != nil {
		return err
	}
	return nil
}
