package docker

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	dc "github.com/fsouza/go-dockerclient"
	"github.com/pkg/errors"
)

// NewClient of docker
func NewClient(endpoint string) (*dc.Client, error) {
	if endpoint == "" {
		if os.Getenv("DOCKER_MACHINE_NAME") != "" {
			client, err := dc.NewClientFromEnv()
			if err != nil {
				return nil, errors.Wrap(err, "failed to create client from environment")
			}
			return client, nil
		} else if os.Getenv("DOCKER_HOST") != "" {
			endpoint = os.Getenv("DOCKER_HOST")
		} else if os.Getenv("DOCKER_URL") != "" {
			endpoint = os.Getenv("DOCKER_URL")
		} else if runtime.GOOS == "windows" {
			endpoint = "http://localhost:2375"
		} else {
			endpoint = "unix:///var/run/docker.sock"
		}
	}

	if os.Getenv("DOCKER_CERT_PATH") != "" && shouldPreferTLS(endpoint) {
		return newTLSClient(endpoint, os.Getenv("DOCKER_CERT_PATH"))
	}

	client, err := dc.NewClient(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return client, nil
}

func shouldPreferTLS(endpoint string) bool {
	return !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "unix://")
}

func newTLSClient(endpoint, certpath string) (*dc.Client, error) {
	ca := fmt.Sprintf("%s/ca.pem", certpath)
	cert := fmt.Sprintf("%s/cert.pem", certpath)
	key := fmt.Sprintf("%s/key.pem", certpath)

	client, err := dc.NewTLSClient(endpoint, cert, key, ca)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return client, nil
}
