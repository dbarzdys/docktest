package docker

import (
	"reflect"

	dc "github.com/fsouza/go-dockerclient"
)

type Remover interface {
	Remove(ids ...string) error
}

type remover struct {
	client *dc.Client
}

func NewRemover(c *dc.Client) Remover {
	return remover{c}
}

func (r remover) Remove(ids ...string) (err error) {
	for _, id := range ids {
		opts := dc.RemoveContainerOptions{
			ID:            id,
			Force:         true,
			RemoveVolumes: true,
		}
		if err = r.client.RemoveContainer(opts); err != nil {
			if reflect.TypeOf(err) == reflect.TypeOf(&dc.NoSuchContainer{}) {
				err = nil
				continue
			}
			return
		}
	}
	return
}
