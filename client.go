package docker

import (
	//"fmt"
	"net/http"

	"golang.org/x/crypto/ssh"

	//"github.com/docker/docker/api/types"
	//"github.com/docker/docker/api/types/filters"
	//"github.com/docker/docker/api/types/swarm"
	docker "github.com/docker/docker/client"
	//"golang.org/x/net/context"
	//"github.com/mitchellh/mapstructure"
)

// Client embeds just docker client
type Client struct {
	*docker.Client
}

// NewSSHClient returns a docker client connected to server using ssh connection
func NewSSHClient(host string, unixSocket string, apiVersion string, sshConfig *ssh.ClientConfig) (*Client, error) {
	var dockerClient *docker.Client

	client := &Client{}

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
func NewClient(host string, apiVersion string, httpClient *http.Client) (*Client, error) {
	var dockerClient *docker.Client

	client := &Client{}

	dockerClient, err := docker.NewClient(host, apiVersion, httpClient, nil)
	if err != nil {
		return nil, err
	}

	client.Client = dockerClient

	return client, nil
}

// NewDefaultClient returns a docker client using default function
func NewDefaultClient() (*Client, error) {

	client := &Client{}

	dockerClient, err := docker.NewEnvClient()
	if err != nil {
		return nil, err
	}

	client.Client = dockerClient

	return client, nil
}
