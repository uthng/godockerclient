package docker

import (
	"fmt"
	"regexp"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
)

// SwarmGetServices returns a list of swarm services created in the cluster
// with options given
func (client *Client) SwarmGetServices(options map[string]string) ([]swarm.Service, error) {
	if len(options) <= 0 {
		// Get swarm services without filter
		return client.ServiceList(client.ctx, types.ServiceListOptions{})
	}

	opts := filters.NewArgs()
	for key, val := range options {
		opts.Add(key, val)
	}

	return client.ServiceList(client.ctx, types.ServiceListOptions{Filters: opts})
}

// SwarmFindServiceByID searches and returns the swarm service corresponding to
// the given ID
func (client *Client) SwarmFindServiceByID(id string, services []swarm.Service) (*swarm.Service, error) {
	for _, service := range services {
		if service.ID == id {
			return &service, nil
		}
	}

	return nil, fmt.Errorf("No service with ID %s found", id)
}

// SwarmFindServiceByName searchs and returns the swarm service corresponding
// to the given pattern
func (client *Client) SwarmFindServiceByName(pattern string, services []swarm.Service) []swarm.Service {
	var srvs []swarm.Service

	for _, service := range services {
		matched, err := regexp.MatchString(pattern, service.Spec.Annotations.Name)
		if err != nil && matched {
			srvs = append(srvs, service)
		}
	}

	return srvs
}
