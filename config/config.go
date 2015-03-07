package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	Reader      string
	Bibfiles    []string
	SearchRoots []string
}

type ConfigNotFoundError struct {
	i int
}

func (e *ConfigNotFoundError) Error() string {
	return fmt.Sprintf("No configuration file found")
}

// Search locations for a config file
// In order:
//	1. /home/
//	2. current directory
//	3. $GOPATH/bin/			[not implemented]
func FindConfig() (string, error) {
	homePath := os.Getenv("HOME")
	curPath, _ := os.Getwd()
	roots := []string{homePath, curPath}

	var files []os.FileInfo
	var err error
	for _, root := range roots {
		files, err = ioutil.ReadDir(root)
		if err != nil {
			fmt.Println(err)
		}
		for _, f := range files {
			if f.Name() == ".peer2.yaml" {
				return filepath.Join(root, f.Name()), nil
			}
		}
	}
	return "", &ConfigNotFoundError{}
}

// Parse a configuration file
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
