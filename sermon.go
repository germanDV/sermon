package sermon

import (
	"net/http"
	"os"
	"sync"

	"gitlab.com/germandv/sermon/internal/httpclient"
	"gitlab.com/germandv/sermon/sermonconfig"
	"gitlab.com/germandv/sermon/sermoncore"
	"gitlab.com/germandv/sermon/sermonreport"
)

// Check verifies the health of a Service.
func Check(s sermoncore.Service) *sermoncore.ServiceStatus {
	client := httpclient.New(&http.Client{Timeout: s.Timeout.Duration})
	err := s.Health(client)
	return &sermoncore.ServiceStatus{
		Name:    s.Name,
		Healthy: err == nil,
		Err:     err,
	}
}

// CheckAll verifies the health of all services listed in the config.
func CheckAll(config *sermonconfig.Config) *sermonreport.Report {
	report := &sermonreport.Report{}
	var wg sync.WaitGroup

	for name, service := range config.Services {
		wg.Add(1)
		s := service
		s.Name = name

		go func() {
			defer wg.Done()
			checkWithRetry := withRetry(config.Attempts.Value, Check, func(ss *sermoncore.ServiceStatus) bool {
				return !ss.Healthy
			})
			report.Add(checkWithRetry(s))
		}()
	}

	wg.Wait()
	return report
}

// Run parses the config, checks all services and emails the results.
func Run(configFileContent string) error {
	config, err := sermonconfig.Parse(configFileContent)
	if err != nil {
		return err
	}

	report := CheckAll(config)
	report.Log(os.Stdout)
	err = report.EmailFail(config.Email.Address)
	if err != nil {
		return err
	}

	return nil
}

// withRetry re-runs a function a given number of times, as long as the
// shouldRetry function returns `true`.
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
