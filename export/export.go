package export

import (
	"fmt"
	"os"
	"strings"

	"github.com/dbarzdys/docktest/docker"
)

func Create(
	path string,
	m map[string]string,
	containers []docker.Container,
) (err error) {
	f, err := os.Create(path)
	if err != nil {
		err = fmt.Errorf("got error while creating export file: %v", err)
		return
	}
	defer f.Close()
	for k, v := range m {
		v = os.Expand(v, func(s string) string {
			original := fmt.Sprintf("${%s}", s)
			path := strings.Split(s, ".")
			if len(path) != 3 {
				return original
			}
			switch path[0] {
			case "services":
				found := -1
				for i, c := range containers {
					if c.Name == path[1] {
						found = i
					}
				}
				if found == -1 && path[2] != "ip" {
					return original
				}
				return containers[found].IP
			default:
				return original
			}
		})
		line := fmt.Sprintf("%s=%s\n", k, v)
		_, err = f.WriteString(line)
		if err != nil {
			return
		}
	}
	return
}

func Exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
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
