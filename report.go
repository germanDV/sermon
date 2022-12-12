package sermon

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/smtp"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"text/template"
)

// ServiceStatus contains information about a service after checking its health.
type ServiceStatus struct {
	Name    string
	Healthy bool
	Err     error
}

// Report consolidates information about health of all services.
type Report struct {
	Services   []*ServiceStatus
	Successful int
	Failed     int
	mu         sync.Mutex
}

// Add adds information about a service to a Report in a concurrency-safe fashion.
func (r *Report) Add(service *ServiceStatus) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if service.Healthy {
		r.Successful++
	} else {
		r.Failed++
	}
	r.Services = append(r.Services, service)
}

// Log prints Report information to stdout.
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

// Email sends Report via email.
func (r *Report) Email(to string) error {
	username, okU := os.LookupEnv("EMAIL_USERNAME")
	password, okP := os.LookupEnv("EMAIL_PASSWORD")
	if !okU || !okP {
		return errors.New("The following env vars must be present to be able to email the report: EMAIL_USERNAME, EMAIL_PASSWORD")
	}

	host := os.Getenv("EMAIL_HOST")
	if host == "" {
		host = "smtp.gmail.com"
	}

	portStr := os.Getenv("EMAIL_PORT")
	if portStr == "" {
		portStr = "587"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return errors.New("EMAIL_PORT must be a number")
	}

	var msg bytes.Buffer
	r.Log(&msg)

	email := struct {
		To   string
		Body string
	}{
		To:   to,
		Body: msg.String(),
	}

	tpl, err := template.ParseFiles(filepath.Join("templates", "email.tmpl"))
	if err != nil {
		return err
	}

	var content bytes.Buffer
	err = tpl.Execute(&content, email)

	auth := smtp.PlainAuth("", username, password, host)

	err = smtp.SendMail(
		fmt.Sprintf("%s:%d", host, port),
		auth,
		username,
		[]string{to},
		content.Bytes(),
	)

	return err
}

// EmailFail sends the Report via email only if there are unhealthy services.
func (r *Report) EmailFail(to string) error {
	someUnhealthy := some(r.Services, func(ss *ServiceStatus) bool {
		return !ss.Healthy
	})

	if someUnhealthy {
		return r.Email(to)
	}

	fmt.Println("All services healthy, skipping email.")
	return nil
}
