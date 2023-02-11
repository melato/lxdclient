package lxdclient

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

func ReadYaml(file string, v interface{}) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, v)
	if err != nil {
		return fmt.Errorf("%s: %v", file, err)
	}
	return nil
}
