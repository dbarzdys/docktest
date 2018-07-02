package docker

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/dbarzdys/docktest/config"
	dc "github.com/fsouza/go-dockerclient"
)

type DepList struct {
	Map   map[string]config.Service
	Order []string
}

func NewDepList(m map[string]config.Service) DepList {
	order := make([]string, 0)
	for k := range m {
		order = append(order, k)
	}
	dl := DepList{
		Map:   m,
		Order: order,
	}
	return dl
}

func (list *DepList) Len() int {
	return len(list.Order)
}

func (list *DepList) Less(i, j int) bool {
	iv := list.Map[list.Order[i]]
	for _, d := range iv.DependsOn {
		if d == list.Order[j] {
			return false
		}
	}
	return true
}

func (list *DepList) Swap(i, j int) {
	tmp := list.Order[i]
	list.Order[i] = list.Order[j]
	list.Order[j] = tmp
}

type Container struct {
	ID   string
	Name string
	IP   string
}

type Runner interface {
	Run(map[string]config.Service) ([]Container, error)
}

type runner struct {
	client *dc.Client
}

// NewRunner creates a new runner
func NewRunner(c *dc.Client) Runner {
	return runner{c}
}

func (r runner) Run(services map[string]config.Service) (
	[]Container,
	error,
) {
	dl := NewDepList(services)
	sort.Sort(&dl)
	containers := make([]Container, 0)
	for _, name := range dl.Order {
		svc := services[name]
		for k, v := range svc.Env {
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
					if found == -1 || path[2] != "ip" {
						return original
					}
					return containers[found].IP
				default:
					return original
				}
			})
			svc.Env[k] = v
		}
		c, err := r.runOne(name, svc)
		if err != nil {
			return nil, err
		}
		containers = append(containers, c)
	}
	return containers, nil
}

func makeRandomName() string {
	r := make([]byte, 32)
	rand.Read(r)
	h := hex.EncodeToString(r)
	return fmt.Sprintf("test_%s", h)
}

func (r runner) runOne(name string, conf config.Service) (
	container Container,
	err error,
) {
	tag := conf.Tag
	image := conf.Image

	if tag == "" {
		tag = "latest"
	}
	env := make([]string, 0)
	for k, v := range conf.Env {
		e := fmt.Sprintf("%s=%s", k, v)
		env = append(env, e)
	}
	repo := fmt.Sprintf("%s:%s", image, tag)
	if _, err = r.client.InspectImage(repo); err != nil {
		err = fmt.Errorf("could not find image %s, error: %v", repo, err)
		return
	}
	c, err := r.client.CreateContainer(dc.CreateContainerOptions{
		Name: makeRandomName(),
		Config: &dc.Config{
			Image: repo,
			Env:   env,
		},
		HostConfig: &dc.HostConfig{
			RestartPolicy: dc.RestartOnFailure(3),
		},
	})

	if err != nil {
		err = fmt.Errorf("could not create container: %v", err)
		return
	}
	if err = r.client.StartContainer(c.ID, nil); err != nil {
		err = fmt.Errorf("could not start container: %v", err)
		return
	}

	c, err = r.client.InspectContainer(c.ID)
	if err != nil {
		err = fmt.Errorf("could not inspect container: %v", err)
		return
	}
	container.IP = c.NetworkSettings.IPAddress
	container.ID = c.ID
	container.Name = name
	return
}
