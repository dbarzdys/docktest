package docker

import (
	"errors"
	"fmt"

	"github.com/dbarzdys/docktest/config"
	dc "github.com/fsouza/go-dockerclient"
)

type Puller interface {
	Pull(map[string]config.Service) error
}

type puller struct {
	client *dc.Client
}

func NewPuller(client *dc.Client) Puller {
	return puller{client}
}

func (p puller) Pull(m map[string]config.Service) error {
	authMap, err := dc.NewAuthConfigurationsFromDockerCfg()
	if err != nil {
		return err
	}
	var auth *dc.AuthConfiguration
	for _, cfg := range authMap.Configs {
		auth = &cfg
		break
	}
	if auth == nil {
		err = errors.New("could not find any docker configuration")
		return err
	}
	for _, svc := range m {
		if svc.Tag == "" {
			svc.Tag = "latest"
		}
		opts := dc.PullImageOptions{
			Repository: svc.Image,
			Tag:        svc.Tag,
		}
		fmt.Printf("pulling: %s:%s\n", svc.Image, svc.Tag)

		err = p.client.PullImage(opts, *auth)
		if err != nil {
			err = fmt.Errorf("got error while pulling image: %v", err)
			return err
		}
	}
	return nil
}
