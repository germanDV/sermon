package main

import (
	"fmt"
	"os"
	"sync"

	"gitlab.com/germandv/sermon"
)

func checkAll(services map[string]sermon.Service) {
	var wg sync.WaitGroup
	for name, service := range services {
		s := service
		wg.Add(1)
		s.Name = name
		go func() {
			err := s.Health(&wg)
			if err != nil {
				log(s.Name, err)
			} else {
				log(s.Name, nil)
			}
		}()
	}
	wg.Wait()
}

func log(serviceName string, err error) {
	if err != nil {
		fmt.Printf("[ERROR] GET %s: %s\n", serviceName, err)
	} else {
		fmt.Printf("[OK] GET %s\n", serviceName)
	}
}

func main() {
	// TODO: embed "services.toml" file
	services, err := sermon.ReadConfig("services.toml")
	if err != nil {
		fmt.Printf("Error reading config file, %s\n", err)
		os.Exit(1)
	}

	checkAll(services)
}
