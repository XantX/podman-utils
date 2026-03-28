package podman

import (
	"os/exec"
	"strings"
)

type Container struct {
	ID   string
	Name string
}

type ContainerDetail struct {
	ID      string
	Name    string
	Image   string
	Status  string
	Ports   string
	Created string
	Command string
	State   string
}

type Client struct{}

func New() *Client {
	return &Client{}
}

func (c *Client) ListStoppedContainers() ([]Container, error) {
	cmd := exec.Command("podman", "ps", "-a", "--filter", "status!=running", "--format", "{{.ID}}|{{.Names}}")
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

func (c *Client) ListRunningContainers() ([]Container, error) {
	cmd := exec.Command("podman", "ps", "--format", "{{.ID}}|{{.Names}}")
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

func (c *Client) ListRunningContainersDetailed() ([]ContainerDetail, error) {
	cmd := exec.Command("podman", "ps", "-a", "--format", "{{.ID}}|{{.Names}}|{{.Image}}|{{.Status}}|{{.Ports}}|{{.CreatedAt}}|{{.Command}}|{{.State}}")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var containers []ContainerDetail
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 8)
		if len(parts) >= 8 {
			containers = append(containers, ContainerDetail{
				ID:      parts[0],
				Name:    parts[1],
				Image:   parts[2],
				Status:  parts[3],
				Ports:   parts[4],
				Created: parts[5],
				Command: parts[6],
				State:   parts[7],
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

func (c *Client) StopContainer(id string) error {
	cmd := exec.Command("podman", "stop", id)
	_, err := cmd.Output()
	return err
}
