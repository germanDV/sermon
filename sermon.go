package sermon

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"sync"
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
		return nil
	}
	return fmt.Errorf("Invalid email address: %s", text)
}

type Service struct {
	Name     string
	Endpoint Endpoint
	Codes    []StatusCode
	Timeout  Timeout
	Email    Email
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

func ReadConfig(config string) (map[string]Service, error) {
	services := map[string]Service{}

	_, err := toml.Decode(config, &services)
	if err != nil {
		return nil, err
	}

	return services, nil
}

func (s *Service) Health(wg *sync.WaitGroup) error {
	defer wg.Done()

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
