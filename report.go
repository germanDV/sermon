package sermon

import (
	"fmt"
	"sync"
)

// ServiceStatus contains information about a service after checking its health.
type ServiceStatus struct {
	Name    string
	Healthy bool
	Err     error
}

// Report consolidates information about health of all services.
type Report struct {
	Services   []*ServiceStatus
	Successful int
	Failed     int
	mu         sync.Mutex
}

// Add adds information about a service to a Report in a concurrency-safe fashion.
func (r *Report) Add(service *ServiceStatus) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if service.Healthy {
		r.Successful++
	} else {
		r.Failed++
	}
	r.Services = append(r.Services, service)
}

// Log prints Report information to stdout.
func (r *Report) Log() {
	fmt.Printf("SUCCESSFUL: %d\n", r.Successful)
	fmt.Printf("FAILED: %d\n", r.Failed)
	fmt.Printf("TOTAL: %d\n\n", r.Successful+r.Failed)
	for _, service := range r.Services {
		if !service.Healthy {
			fmt.Printf("GET %s -> ERROR: %s\n", service.Name, service.Err)
		} else {
			fmt.Printf("GET %s -> OK\n", service.Name)
		}
	}
}

// Email sends Report via email.
// TODO: implement!
func (r *Report) Email(to string) {
	fmt.Printf("Emailing report to %q (WIP).\n", to)
}
