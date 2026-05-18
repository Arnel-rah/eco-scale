package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type DockerScanner struct {
	cli *client.Client
}

func (ds *DockerScanner) StopContainer(d string) any {
	panic("unimplemented")
}

func NewDockerScanner() (*DockerScanner, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("impossible de se connecter à Docker: %w", err)
	}

	return &DockerScanner{cli: cli}, nil
}

func (ds *DockerScanner) ListActiveContainers() ([]types.Container, error) {
	ctx := context.Background()

	containers, err := ds.cli.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("impossible de lister les conteneurs: %w", err)
	}

	return containers, nil
}

func (ds *DockerScanner) Close() {
	if ds.cli != nil {
		ds.cli.Close()
	}
}
