package podman

import (
	"os/exec"
	"strings"
)

type Container struct {
	ID   string
	Name string
}

type Client struct{}

func New() *Client {
	return &Client{}
}

func (c *Client) ListStoppedContainers() ([]Container, error) {
	cmd := exec.Command("podman", "ps", "-a", "--filter", "status=exited", "--format", "{{.ID}}|{{.Names}}")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var containers []Container
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 2)
		if len(parts) == 2 {
			containers = append(containers, Container{
				ID:   parts[0],
				Name: parts[1],
			})
		}
	}
	return containers, nil
}

func (c *Client) StartContainer(id string) error {
	cmd := exec.Command("podman", "start", id)
	_, err := cmd.Output()
	return err
}
