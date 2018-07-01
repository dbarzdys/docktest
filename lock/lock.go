package lock

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type LockConfig struct {
	ExportFile string   `yaml:"export_file"`
	Containers []string `yaml:"containers"`
}

func Create(path string, conf LockConfig) (err error) {
	f, err := os.Create(path)
	if err != nil {
		err = fmt.Errorf("got error while creating lock file: %v", err)
		return
	}
	defer f.Close()
	encoder := yaml.NewEncoder(f)
	err = encoder.Encode(&conf)
	if err != nil {
		err = fmt.Errorf("got error while encoding lock file: %v", err)
		return
	}
	defer encoder.Close()
	return
}

func Exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func Read(path string) (
	conf LockConfig,
	err error,
) {
	if !Exists(path) {
		err = fmt.Errorf("could not find lock file at: %s", path)
		return
	}
	f, err := os.Open(path)
	if err != nil {
		err = fmt.Errorf("got error while trying to open lock file: %v", err)
		return
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&conf)
	if err != nil {
		err = fmt.Errorf("got error while trying to decode lock file: %v", err)
		return
	}
	return
}

func Remove(path string) (err error) {
	if !Exists(path) {
		err = fmt.Errorf("could not find lock file at: %s", path)
		return
	}
	err = os.Remove(path)
	if err != nil {
		err = fmt.Errorf("could not remove lock file: %v", err)
	}
	return
}
