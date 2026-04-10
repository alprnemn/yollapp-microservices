package consul

import (
	"context"
	"errors"
	"fmt"
	"github.com/alprnemn/yollapp-microservices/shared/discovery"
	consul "github.com/hashicorp/consul/api"
	"strconv"
	"strings"
)

type Registry struct {
	client *consul.Client
}

func New(addr string) (*Registry, error) {
	config := consul.DefaultConfig()
	config.Address = addr
	cl, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &Registry{client: cl}, nil
}

func (r *Registry) Register(ctx context.Context, instanceID string, serviceName string, serviceAddress string) error {

	parts := strings.Split(serviceAddress, ":")
	if len(parts) != 2 {
		return errors.New("service address format must be <host>:<port>")
	}

	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}

	agentServiceRegistration := &consul.AgentServiceRegistration{
		Address: parts[0],
		ID:      instanceID,
		Name:    serviceName,
		Port:    port,
		Check: &consul.AgentServiceCheck{
			CheckID: instanceID,
			TTL:     "5s",
		},
	}
	return r.client.Agent().ServiceRegister(agentServiceRegistration)
}

func (r *Registry) Deregister(ctx context.Context, instanceID string, _ string) error {
	return r.client.Agent().ServiceDeregister(instanceID)
}

func (r *Registry) ServiceAddresses(ctx context.Context, serviceName string) ([]string, error) {
	entries, _, err := r.client.Health().Service(
		serviceName,
		"",
		true,
		nil,
	)
	if err != nil {
		return nil, err
	} else if len(entries) == 0 {
		return nil, discovery.ErrNotFound
	}
	var res []string
	for _, e := range entries {
		res = append(res, fmt.Sprintf("%s:%d", e.Service.Address,
			e.Service.Port))
	}
	return res, nil
}

func (r *Registry) ReportHealthyState(instanceID string, _ string) error {
	return r.client.Agent().PassTTL(instanceID, "")
}
