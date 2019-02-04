package docker

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/docker/docker/api/types"
	//"github.com/docker/docker/api/types/swarm"

	"github.com/uthng/common/ssh"
)

var client *Client
var ctx = context.Background()

func TestNewSSHClient(t *testing.T) {
	sshUser := os.Getenv("SSHUSER")
	sshKey := os.Getenv("SSHKEY")
	sshHost := os.Getenv("SSHHOST")

	config, err := ssh.NewClientConfigWithKeyFile(sshUser, sshKey, "", 0, false)
	assert.Nil(t, err)

	client, err = NewSSHClient(ctx, sshHost, "/var/run/docker.sock", "1.30", config.ClientConfig)
	assert.Nil(t, err)

	services, err := client.ServiceList(client.ctx, types.ServiceListOptions{})
	assert.Nil(t, err)

	for _, service := range services {
		fmt.Println(service.ID, " ", service.Spec.Annotations.Name)
	}

}
