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

// SwarmFindServiceByName searchs and returns the swarm services
// that their names match the given pattern
func (client *Client) SwarmFindServiceByName(pattern string, services []swarm.Service) ([]swarm.Service, error) {
	var srvs []swarm.Service

	for _, service := range services {
		matched, err := regexp.MatchString(pattern, service.Spec.Annotations.Name)
		if err != nil {
			return srvs, err
		}

		if matched {
			srvs = append(srvs, service)
		}
	}

	return srvs, nil
}

// SwarmDeleteServices removes all swarm services with the name matching
// the given pattern. It returns a list of deleted services.
func (client *Client) SwarmDeleteServices(pattern string) ([]swarm.Service, error) {
	var srvs []swarm.Service

	// Get all services
	services, err := client.SwarmGetServices(nil)
	if err != nil {
		return srvs, err
	}

	// Find services matching pattern in its name
	servicesMatched, err := client.SwarmFindServiceByName(pattern, services)
	if err != nil {
		return srvs, err
	}

	// Loop to remove
	for _, s := range servicesMatched {
		matched, err := regexp.MatchString(pattern, s.Spec.Annotations.Name)
		if err != nil {
			return srvs, err
		}

		if matched {
			err := client.ServiceRemove(client.ctx, s.ID)
			if err != nil {
				return srvs, err
			}

			srvs = append(srvs, s)
		}
	}

	return srvs, nil
}
