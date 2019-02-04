package docker

import (
	//"fmt"
	"context"
	"net/http"

	"golang.org/x/crypto/ssh"

	//"github.com/docker/docker/api/types"
	//"github.com/docker/docker/api/types/filters"
	//"github.com/docker/docker/api/types/swarm"
	docker "github.com/docker/docker/client"
	//"github.com/mitchellh/mapstructure"
)

// Client embeds just docker client
type Client struct {
	ctx context.Context

	*docker.Client
}

// NewSSHClient returns a docker client connected to server using ssh connection
func NewSSHClient(ctx context.Context, host string, unixSocket string, apiVersion string, sshConfig *ssh.ClientConfig) (*Client, error) {
	var dockerClient *docker.Client

	client := &Client{
		ctx: ctx,
	}

	dialer := &dialerSSH{
		host:   host,
		socket: unixSocket,
		config: sshConfig,
	}
	httpClient := &http.Client{
		Transport: &http.Transport{
			Dial: dialer.Dial,
		},
	}

	dockerClient, err := docker.NewClient("unix://"+unixSocket, apiVersion, httpClient, nil)

	// New docker API moby
	//dockerClient, err := docker.NewClientWithOpts(docker.WithHost("unix://"+unixSocket), docker.WithVersion(apiVersion), docker.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}

	client.Client = dockerClient

	return client, nil
}

// NewClient returns a docker client with given paramters
func NewClient(ctx context.Context, host string, apiVersion string, httpClient *http.Client) (*Client, error) {
	var dockerClient *docker.Client

	client := &Client{
		ctx: ctx,
	}

	dockerClient, err := docker.NewClient(host, apiVersion, httpClient, nil)
	if err != nil {
		return nil, err
	}

	client.Client = dockerClient

	return client, nil
}

// NewDefaultClient returns a docker client using default function
func NewDefaultClient(ctx context.Context) (*Client, error) {

	client := &Client{
		ctx: ctx,
	}

	dockerClient, err := docker.NewEnvClient()
	if err != nil {
		return nil, err
	}

	client.Client = dockerClient

	return client, nil
}
