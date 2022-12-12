package sermon

import (
	"fmt"
	"sync"
)

// Health makes an HTTP request to check the health of the service.
func (s *Service) Health() error {
	status, err := get(s.Endpoint, s.Timeout)
	if err != nil {
		return err
	}

	if !in(s.Codes, status) {
		e := fmt.Errorf("Got status %d, want one of %v", status, s.Codes)
		return e
	}

	return nil
}

// Check verifies the health of a Service.
func Check(s Service) *ServiceStatus {
	err := s.Health()
	return &ServiceStatus{
		Name:    s.Name,
		Healthy: err == nil,
		Err:     err,
	}
}

// CheckAll verifies the health of all services listed in the config.
func CheckAll(config *Config) *Report {
	report := &Report{}
	var wg sync.WaitGroup

	for name, service := range config.Services {
		wg.Add(1)
		s := service
		s.Name = name

		go func() {
			defer wg.Done()
			checkWithRetry := withRetry(config.Attempts.Value, Check, func(ss *ServiceStatus) bool {
				return !ss.Healthy
			})
			report.Add(checkWithRetry(s))
		}()
	}

	wg.Wait()
	return report
}

// Run parses the config, checks all services and emails the results.
func Run(configFileContent string) error {
	config, err := ReadConfig(configFileContent)
	if err != nil {
		return err
	}

	report := CheckAll(config)
	report.Log()

	err = report.Email(config.Email.Address)
	if err != nil {
		return err
	}

	return nil
}
