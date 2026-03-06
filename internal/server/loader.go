package server

import (
	"os"

	"gopkg.in/yaml.v3"
)

func LoadServer(path string) (*Server, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var s Server

	err = yaml.Unmarshal(data, &s)

	if err != nil {
		return nil, err
	}

	return &s, nil
}
