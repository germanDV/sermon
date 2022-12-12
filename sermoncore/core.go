package sermoncore

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Timeout struct {
	Duration time.Duration
}

func (t *Timeout) UnmarshalText(text []byte) error {
	var err error
	t.Duration, err = time.ParseDuration(string(text))
	return err
}

type Endpoint struct {
	URL *url.URL
}

func (e *Endpoint) UnmarshalText(text []byte) error {
	var err error
	e.URL, err = url.ParseRequestURI(string(text))
	return err
}

type StatusCode struct {
	Code int
}

func (s *StatusCode) UnmarshalText(text []byte) error {
	code, err := strconv.Atoi(string(text))
	if err != nil || code < 100 || code > 599 {
		return fmt.Errorf("Invalid status code: %s", text)
	}
	s.Code = code
	return nil
}

// Service represents a web service which health is to be monitored.
type Service struct {
	Name     string
	Endpoint Endpoint
	Codes    []StatusCode
	Timeout  Timeout
}

// ServiceStatus contains information about a service after checking its health.
type ServiceStatus struct {
	Name    string
	Healthy bool
	Err     error
}

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

// in checks if the given item is included in the given slice of items.
func in(items []StatusCode, item int) bool {
	for _, i := range items {
		if item == i.Code {
			return true
		}
	}
	return false
}

// get makes a GET HTTP request and returns the response status code.
func get(endpoint Endpoint, timeout Timeout) (int, error) {
	client := &http.Client{Timeout: timeout.Duration}
	resp, err := client.Get(endpoint.URL.String())
	if err != nil {
		return 0, err
	}
	return resp.StatusCode, nil
}
