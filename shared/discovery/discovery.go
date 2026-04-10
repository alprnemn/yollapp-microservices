package discovery

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Registry defines a service registry
type Registry interface {
	// Register creates a service instance record
	// in the registry
	Register(ctx context.Context, instanceID string, serviceName string, serviceAddress string) error

	// Deregister removes a service instance record
	// from registry
	Deregister(ctx context.Context, instanceID string, serviceName string) error

	// ServiceAddresses returns the list of addresses of active instances
	// of the given service
	ServiceAddresses(ctx context.Context, serviceName string) ([]string, error)

	// ReportHealthyState is a push mechanism for reporting
	// healthy state to the registry
	ReportHealthyState(instanceID string, serviceName string) error
}

var ErrNotFound = errors.New("no service addresses found")

// GenerateInstanceID generates a psuedo-random service
// instance identifier, using a service name
// suffixed by dash and a random number
func GenerateInstanceID(serviceName string) string {
	return fmt.Sprintf("%s-%d", serviceName, rand.New(rand.NewSource(time.Now().UnixNano())).Int())
}
