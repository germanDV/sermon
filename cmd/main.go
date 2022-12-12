package main

import (
	_ "embed"
	"fmt"
	"os"
	"sync"

	"gitlab.com/germandv/sermon"
)

func checkAll(config *sermon.Config) *sermon.Report {
	summary := &sermon.Report{}
	var wg sync.WaitGroup

	for name, service := range config.Services {
		s := service
		wg.Add(1)
		s.Name = name
		go func() {
			defer wg.Done()
			err := s.Health()
			summary.Add(&sermon.ServiceStatus{
				Name:    s.Name,
				Healthy: err == nil,
				Err:     err,
			})
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
