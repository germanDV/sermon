package sermon

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

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

type Email struct {
	Address string
}

func (e *Email) UnmarshalText(text []byte) error {
	if EmailRX.Match(text) {
		e.Address = string(text)
		return nil
	}
	return fmt.Errorf("Invalid email address: %s", text)
}

type Attempts struct {
	Value int
}

func (a *Attempts) UnmarshalText(text []byte) error {
	n, err := strconv.Atoi(string(text))
	if err != nil || n < 1 || n > 10 {
		return fmt.Errorf("Invalid number of attempts (min 1, max 10): %s", text)
	}
	a.Value = n
	return nil
}

// Service represents a web service which health is to be monitored.
type Service struct {
	Name     string
	Endpoint Endpoint
	Codes    []StatusCode
	Timeout  Timeout
}

// Config represents the structure of the TOML file that lists the services
// to be checked and some common settings.
type Config struct {
	Email    Email
	Attempts Attempts
	Services map[string]Service
}

func in(items []StatusCode, item int) bool {
	for _, i := range items {
		if item == i.Code {
			return true
		}
	}
	return false
}

func get(endpoint Endpoint, timeout Timeout) (int, error) {
	client := &http.Client{Timeout: timeout.Duration}
	resp, err := client.Get(endpoint.URL.String())
	if err != nil {
		return 0, err
	}
	return resp.StatusCode, nil
}

// ReadConfig parses the TOML file that lists the services to monitor.
func ReadConfig(config string) (*Config, error) {
	cfg := &Config{}

	_, err := toml.Decode(config, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
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
