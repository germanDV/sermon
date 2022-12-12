package main

import (
	_ "embed"
	"fmt"
	"os"
	"sync"

	"gitlab.com/germandv/sermon"
)

func withRetry[T any, U any](
	maxAttempts int,
	fn func(service T) *U,
	shouldRetry func(status *U) bool,
) func(item T) *U {
	attempts := 0

	return func(item T) *U {
		result := new(U)
		for attempts < maxAttempts {
			attempts++
			result = fn(item)
			if !shouldRetry(result) {
				break
			}
		}
		return result
	}
}

func check(service sermon.Service) *sermon.ServiceStatus {
	err := service.Health()
	return &sermon.ServiceStatus{
		Name:    service.Name,
		Healthy: err == nil,
		Err:     err,
	}
}

func checkAll(config *sermon.Config) *sermon.Report {
	summary := &sermon.Report{}
	var wg sync.WaitGroup

	for name, service := range config.Services {
		wg.Add(1)
		s := service
		s.Name = name

		go func() {
			defer wg.Done()

			checkWithRetry := withRetry(config.Attempts.Value, check, func(ss *sermon.ServiceStatus) bool {
				return !ss.Healthy
			})

			status := checkWithRetry(s)
			summary.Add(status)
		}()
	}

	wg.Wait()
	return summary
}

//go:embed services.toml
var configFile string

func main() {
	config, err := sermon.ReadConfig(configFile)
	if err != nil {
		fmt.Printf("Error reading config file, %s\n", err)
		os.Exit(1)
	}

	summary := checkAll(config)
	summary.Log()
	summary.Email(config.Email.Address)
}
