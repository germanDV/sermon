package main

import (
	_ "embed"
	"fmt"
	"os"
	"sync"

	"gitlab.com/germandv/sermon"
	"gitlab.com/germandv/sermon/internal/report"
)

func checkAll(config *sermon.Config) *report.Report {
	summary := &report.Report{}
	var wg sync.WaitGroup

	for name, service := range config.Services {
		s := service
		wg.Add(1)
		s.Name = name
		go func() {
			defer wg.Done()
			err := s.Health()
			summary.Add(&report.ServiceStatus{
				Name:    s.Name,
				Healthy: err == nil,
				Err:     err,
			})
		}()
	}

	wg.Wait()
	return summary
}

func log(serviceName string, err error) {
	if err != nil {
		fmt.Printf("[ERROR] GET %s: %s\n", serviceName, err)
	} else {
		fmt.Printf("[OK] GET %s\n", serviceName)
	}
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
