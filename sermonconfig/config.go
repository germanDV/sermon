package sermonconfig

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
	"gitlab.com/germandv/sermon/sermoncore"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

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

// Config represents the structure of the TOML file that lists the services
// to be checked and some common settings.
type Config struct {
	Email    Email
	Attempts Attempts
	Services map[string]sermoncore.Service
}

// Parse parses the TOML file that lists the services to monitor.
func Parse(config string) (*Config, error) {
	cfg := &Config{}

	_, err := toml.Decode(config, cfg)
	if err != nil {
		return nil, err
	}

	if cfg.Email.Address == "" {
		return nil, errors.New("Missing `email`")
	}
	if cfg.Attempts.Value == 0 {
		return nil, errors.New("Missing `attempts`")
	}

	for name, s := range cfg.Services {
		if s.Endpoint.URL == nil {
			return nil, fmt.Errorf("Missing `endpoint` for service %s", name)
		}
		if len(s.Codes) == 0 {
			return nil, fmt.Errorf("Missing `codes` for service %s", name)
		}
		if s.Timeout.Duration == time.Duration(0) {
			return nil, fmt.Errorf("Missing `timeout` for service %s", name)
		}
	}

	return cfg, nil
}
