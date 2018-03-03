package docker

import (
    "fmt"
    "net/http"

    "github.com/pkg/errors"

    "golang.org/x/crypto/ssh"

    "golang.org/x/net/context"
    docker "github.com/docker/docker/client"
    "github.com/docker/docker/api/types"
    "github.com/docker/docker/api/types/swarm"
    "github.com/docker/docker/api/types/filters"
)

var (
    ErrNetworkIDNotFound = errors.New("Network ID not found")
    ErrServiceIDNotFound = errors.New("Service ID not found")
    ErrServiceNameNotFound = errors.New("Service Name not found")
    ErrNodeIDNotFound = errors.New("Node ID not found")
)

type Client struct {
    *docker.Client
}

// NewSSHClient returns a docker client connected to server using ssh connection
func NewSSHClient(host string, unixSocket string, apiVersion string, sshConfig *ssh.ClientConfig) (*Client, error) {
    var dockerClient *docker.Client

    client := &Client{}

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

    dockerClient, err := docker.NewClient("unix://" + unixSocket, apiVersion, httpClient, nil)
    // New version of moby client
    //newClient, err := client.NewClientWithOpts(client.WithHost("unix://" + unixSocket), client.WithVersion(apiVersion), client.WithHTTPClient(httpClient))
    if err != nil {
        return nil, err
    }

    client.Client = dockerClient

    return client, nil
}

// GetSwarmServices returns a list of swarm services created in the cluster
func (client *Client) GetSwarmServices (ctx context.Context, options map[string]string) ([]swarm.Service, error) {
    if len(options) <= 0 {
        // Get swarm services without filter
        return client.ServiceList(ctx, types.ServiceListOptions{})
    } else {
        opts := filters.NewArgs()
        for key, val := range options {
            opts.Add(key, val)
        }

        return client.ServiceList(ctx, types.ServiceListOptions{ Filters: opts, })
    }
}

// FindServiceByID searchs and returns the swarm service corresponding to
// the given ID
func (client *Client) FindSwarmServiceByID (id string, services []swarm.Service) (*swarm.Service, error) {
    for _, service := range services {
        if service.ID == id {
            return &service, nil
        }
    }

    return nil, ErrServiceIDNotFound
}

// FindServiceByName searchs and returns the swarm service corresponding to
// the given name
func (client *Client) FindSwarmServiceByName (name string, services []swarm.Service) (*swarm.Service, error) {
    for _, service := range services {
        if service.Spec.Annotations.Name == name {
            return &service, nil
        }
    }

    return nil, ErrServiceNameNotFound
}

// GetNetworks return a list of network defined in the cluster
func (client *Client) GetNetworks (ctx context.Context, options map[string]string) ([]types.NetworkResource, error) {
    if len(options) <= 0 {
        // Get network list without filter
        return client.NetworkList(ctx, types.NetworkListOptions{})
    } else {
        opts := filters.NewArgs()
        for key, val := range options {
            opts.Add(key, val)
        }

        return client.NetworkList(ctx, types.NetworkListOptions{ Filters: opts, })
    }
}

// FindNetworkByID searchs and returns the NetworkResource corresponding to
// the given ID
func (client *Client) FindNetworkByID (id string, networks []types.NetworkResource) (*types.NetworkResource, error) {
    for _, net := range networks {
        if net.ID == id {
            return &net, nil
        }
    }

    return nil, ErrNetworkIDNotFound
}

// GetSwarmTasks return a list of swarm containers in the cluster
func (client *Client) GetSwarmTasks (ctx context.Context, options map[string]string) ([]swarm.Task, error) {
    if len(options) <= 0 {
        // Get task list without filter
        return client.TaskList(ctx, types.TaskListOptions{})
    } else {
        opts := filters.NewArgs()
        for key, val := range options {
            opts.Add(key, val)
        }

        return client.TaskList(ctx, types.TaskListOptions{ Filters: opts, })
    }

}

// FindSwarmTasksByServiceID searchs and returns the tasks
// related to the given service ID
func (client *Client) FindSwarmTasksByServiceID (serviceId string, tasks []swarm.Task) ([]swarm.Task) {
    var result []swarm.Task

    for _, task := range tasks {
        if task.ServiceID == serviceId {
            result = append(result, task)
        }
    }

    return result
}

// GetSwarmNodes return a list of swarm nodes in the cluster
func (client *Client) GetSwarmNodes (ctx context.Context, options map[string]string) ([]swarm.Node, error) {
    if len(options) <= 0 {
        // Get node list without filter
        return client.NodeList(ctx, types.NodeListOptions{})
    } else {
        opts := filters.NewArgs()
        for key, val := range options {
            opts.Add(key, val)
        }

        return client.NodeList(ctx, types.NodeListOptions{ Filters: opts, })
    }

}

// FindNodeByServiceID returns the node corresponding to the given ID
func (client *Client) FindSwarmNodeByID (id string, nodes []swarm.Node) (*swarm.Node, error) {
    for _, node := range nodes {
        if node.ID == id {
            return &node, nil
        }
    }

    return nil, ErrNodeIDNotFound
}

// GetContainers return a list of containers running on
// the current host of the cluster (like docker ps)
//
// For ContainerListOptions, only Filter is supported. Some other options
// are missing
func (client *Client) GetContainers (ctx context.Context, options map[string]string) ([]types.Container, error) {
    if len(options) <= 0 {
        // Get container list without filter
        return client.ContainerList(ctx, types.ContainerListOptions{})
    } else {
        opts := filters.NewArgs()
        for key, val := range options {
            opts.Add(key, val)
        }

        return client.ContainerList(ctx, types.ContainerListOptions{ Filters: opts, })
    }
}

// CreateExec creates an exec instance and return exec id
//func CreateExec (client *client.Client, id string, cmd []string) (*string, context.Context, error) {
    //config := types.ExecConfig {
        ////AttachStdout: true,
        ////AttachStderr: true,
        ////Tty: false,
        ////Detach: false,
        ////DetachKeys: "ctrl-p,ctrl-q",
        //Cmd: cmd,
    //}
    //// Create a exec instance
    //ctx := context.Background()
    //res, err := client.ContainerExecCreate(context.Background(), id, config)
    //if err != nil {
        //return nil, ctx, err
    //}

    //fmt.Println("exec id ", res.ID)
    //return &res.ID, ctx, nil
//}

// StartExec start exec process and attach it to a reader.
// It returns a bufio reader for command output
//
// This function takes temporairly cmd in argument but
// need to be removed in the next release of docker
// because ContainerExecAttach does not take types.ExecConfig anymore
// Instead, it takes types.ExecStartCheck
//func StartExec (client *client.Client, ctx context.Context, execId string, cmd []string) (*bufio.Reader, error) {
    //config := types.ExecConfig {
        //AttachStdout: true,
        ////AttachStderr: true,
        //Tty: false,
        //Detach: false,
        ////DetachKeys: "ctrl-p,ctrl-q",
        ////Cmd: cmd,
    //}

    //// Create a exec instance
    ////err := client.ContainerExecStart(ctx, execId, types.ExecStartCheck{Detach: false, Tty: false})
    ////if err != nil {
        ////fmt.Println("error start exec")
        ////return nil, err
    ////}

    //// Attach an exec
    //res, err := client.ContainerExecAttach(ctx, execId, config)
    //if err != nil {
        //fmt.Println("error attach exec")
        //return nil, err
    //}
    //defer res.Close()

    //line, _, err := res.Reader.ReadLine()
    //fmt.Println("line ", string(line))
    //return res.Reader, nil
//}

func (client *Client) ExecCommand(ctx context.Context, cid string, cmd []string) ([]byte, error) {
    id, err := client.ContainerExecCreate(ctx, cid,
                types.ExecConfig{
                    //WorkingDir:   "/tmp",
                    //Env:          strslice.StrSlice([]string{"FOO=BAR"}),
                    AttachStdout: true,
                    Cmd: cmd,
                    //Cmd:          strslice.StrSlice([]string{"sh", "-c", "env"}),
                })
    if err != nil {
        fmt.Println("error create exec ", err)
        return nil, err
    }

    fmt.Println("id ", id)

    //insp, err := client.ContainerExecInspect(ctx, id.ID)
    //if err != nil {
         //fmt.Println("error inspect exec ", err)
        //return nil, err
    //}

    //err = client.ContainerExecStart(ctx, id.ID, types.ExecStartCheck{
                    //Detach: false,
                    //Tty:    true,
                //})
    //if err != nil {
         //fmt.Println("error start exec ", err)
        //return nil, err
    //}

    //fmt.Println("inspect ", insp)
    //fmt.Println("cid ", cid)
    resp, err := client.ContainerExecAttach(ctx, id.ID,
                types.ExecStartCheck{
                    Detach: false,
                    Tty:    false,
                })
    if err != nil {
        fmt.Println("error attach exec ", err)
        return nil, err
    }
    defer resp.Close()

    fmt.Println("response ", resp)
    line, _, err := resp.Reader.ReadLine()
    fmt.Println("response ", string(line), err)


    //r, err := ioutil.ReadAll(resp.Reader)
    //if err != nil {
        //fmt.Println("error readall ", err)
        //return nil, err
    //}

    return line, err
}

