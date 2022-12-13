package sermonreport

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"text/template"

	"gitlab.com/germandv/sermon/internal/mailer"
	"gitlab.com/germandv/sermon/sermoncore"
)

const (
	DefaultHost = "smtp.gmail.com"
	DefaultPort = 587
)

// Report consolidates information about health of all services.
type Report struct {
	Services   []*sermoncore.ServiceStatus
	Successful int
	Failed     int
	mu         sync.Mutex
}

// Add adds information about a service to a Report in a concurrency-safe fashion.
func (r *Report) Add(service *sermoncore.ServiceStatus) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if service.Healthy {
		r.Successful++
	} else {
		r.Failed++
	}
	r.Services = append(r.Services, service)
}

// Log prints Report information to the given io.Writer.
func (r *Report) Log(w io.Writer) {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("SUCCESSFUL: %d\n", r.Successful))
	sb.WriteString(fmt.Sprintf("FAILED: %d\n", r.Failed))
	sb.WriteString(fmt.Sprintf("TOTAL: %d\n\n", r.Successful+r.Failed))

	for _, service := range r.Services {
		if !service.Healthy {
			sb.WriteString(fmt.Sprintf("GET %s -> ERROR: %s\n", service.Name, service.Err))
		} else {
			sb.WriteString(fmt.Sprintf("GET %s -> OK\n", service.Name))
		}
	}

	fmt.Fprint(w, sb.String())
}

// getEmailConfig creates a mailer.Config with information from env vars.
func getEmailConfig() (*mailer.Config, error) {
	username, okU := os.LookupEnv("EMAIL_USERNAME")
	password, okP := os.LookupEnv("EMAIL_PASSWORD")
	if !okU || !okP {
		return nil, errors.New("The following env vars must be present to be able to email the report: EMAIL_USERNAME, EMAIL_PASSWORD")
	}

	host := os.Getenv("EMAIL_HOST")
	if host == "" {
		host = DefaultHost
	}

	var port int
	portStr := os.Getenv("EMAIL_PORT")
	if portStr != "" {
		var err error
		port, err = strconv.Atoi(portStr)
		if err != nil {
			return nil, errors.New("EMAIL_PORT must be a number")
		}
	} else {
		port = DefaultPort
	}

	return &mailer.Config{
		Username: username,
		Password: password,
		Host:     host,
		Port:     port,
	}, nil
}

// getEmail parses the email template to populate a proper *mailer.Mail.
func getEmail(to string, msg string) (*mailer.Mail, error) {
	emailData := struct {
		To   string
		Body string
	}{
		To:   to,
		Body: msg,
	}

	tpl, err := template.ParseFiles(filepath.Join("templates", "email.tmpl"))
	if err != nil {
		return nil, err
	}

	var content bytes.Buffer
	err = tpl.Execute(&content, emailData)
	if err != nil {
		return nil, err
	}

	return &mailer.Mail{
		To:   []string{to},
		Body: content.Bytes(),
	}, nil
}

// Email sends Report via email.
func (r *Report) Email(to string) error {
	cfg, err := getEmailConfig()
	if err != nil {
		return err
	}

	var msg bytes.Buffer
	r.Log(&msg)

	email, err := getEmail(to, msg.String())
	if err != nil {
		return err
	}

	return mailer.Send(cfg, email)
}

// EmailFail sends the Report via email only if there are unhealthy services.
func (r *Report) EmailFail(to string) error {
	someUnhealthy := some(r.Services, func(ss *sermoncore.ServiceStatus) bool {
		return !ss.Healthy
	})

	if someUnhealthy {
		return r.Email(to)
	}

	fmt.Println("All services healthy, skipping email.")
	return nil
}

// Some applies the given function to every element in the slice and returns
// `true` if at least one of the invocations returned `true`.
func some[T any](arr []T, fn func(T) bool) bool {
	for _, i := range arr {
		if fn(i) {
			return true
		}
	}
	return false
}
