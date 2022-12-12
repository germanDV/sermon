package sermon

import (
	"fmt"
	"sync"
)

type ServiceStatus struct {
	Name    string
	Healthy bool
	Err     error
}

type Report struct {
	Services   []*ServiceStatus
	Successful int
	Failed     int
	mu         sync.Mutex
}

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

func (r *Report) Log() {
	fmt.Printf("SUCCESSFUL: %d\n", r.Successful)
	fmt.Printf("FAILED: %d\n", r.Failed)
	fmt.Printf("TOTAL: %d\n\n", r.Successful+r.Failed)

	for _, service := range r.Services {
		if !service.Healthy {
			fmt.Printf("[ERROR] GET %s: %s\n", service.Name, service.Err)
		} else {
			fmt.Printf("[OK] GET %s\n", service.Name)
		}
	}
}

func (r *Report) Email(to string) {
	fmt.Printf("Emailing report to %q (WIP).\n", to)
}
