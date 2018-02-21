package docker

import (
    "net/http"

    "golang.org/x/crypto/ssh"
    "github.com/docker/docker/client"
)

func NewSSHClient(host string, unixSocket string, apiVersion string, sshConfig *ssh.ClientConfig) (*client.Client, error) {
    dialer := &dialerSSH{
        host: host,
        socket: unixSocket,
        config: sshConfig,
    }
    httpClient := &http.Client{
        Transport: &http.Transport{
            Dial: dialer.Dial,
        },
    }

    newClient, err := client.NewClient("unix://" + unixSocket, apiVersion, httpClient, nil)
    // New version of moby client
    //newClient, err := client.NewClientWithOpts(client.WithHost("unix://" + unixSocket), client.WithVersion(apiVersion), client.WithHTTPClient(httpClient))
    if err != nil {
        return nil, err
    }

    return newClient, nil
}
