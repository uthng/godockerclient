package docker

import (
    "net"
    "fmt"

    "golang.org/x/crypto/ssh"
)

type dialerSSH struct {
    host        string
    socket      string
    config      *ssh.ClientConfig
}

func (d *dialerSSH) Dial(network, addr string) (net.Conn, error) {
    sshAddr := d.host

    // Establish connection with SSH server
    client, err := ssh.Dial("tcp", sshAddr, d.config)
    if err != nil {
        return nil, err
    }

    conn, err := client.Dial("unix", d.socket)
    if err != nil {
        fmt.Println("ssh err ", err)
        return nil, err
    }

    return conn, err
}
