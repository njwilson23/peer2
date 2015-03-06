package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Reader      string
	Bibfiles    []string
	SearchRoots []string
}

func ParseConfig(fnm string) Config {
	var config Config
	configData, err := ioutil.ReadFile(fnm)
	if err != nil {
		fmt.Println("ERROR:", err)
	}
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		fmt.Println("ERROR:", err)
	}
	return config
}
