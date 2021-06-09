package config

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
)

type Web struct {
	Port int `json:"port"`
}

func getWeb() (Web, error) {

	server := Web{}

	data, err := instance.ReadFile("file/web.json")
	if err != nil && err != io.EOF {
		return server, err
	}

	err = json.Unmarshal(data, &server)
	return server, err

}

func (w Web) Path() string {
	return fmt.Sprintf(":%d", w.Port)
}
