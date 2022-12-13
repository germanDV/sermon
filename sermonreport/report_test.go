package sermonreport

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"testing"

	"gitlab.com/germandv/sermon/expect"
	"gitlab.com/germandv/sermon/sermoncore"
)

func TestGetEmail(t *testing.T) {
	t.Run("MissingRequiredEnvVars", func(t *testing.T) {
		config, err := getEmailConfig()
		expect.Nil(t, config)
		expect.Contains(t, err.Error(), "env vars must be present")
	})

	t.Run("UsesEnvVarsToPopulateConfig", func(t *testing.T) {
		username := "some@email.io"
		password := "abc1234"
		host := "smtp.fastmail.com"
		port := 486
		os.Setenv("EMAIL_USERNAME", username)
		os.Setenv("EMAIL_PASSWORD", password)
		os.Setenv("EMAIL_HOST", host)
		os.Setenv("EMAIL_PORT", fmt.Sprint(port))
		config, err := getEmailConfig()
		expect.NoError(t, err)
		expect.Equal(t, config.Username, username)
		expect.Equal(t, config.Password.Expose(), password)
		expect.Equal(t, config.Host, host)
		expect.Equal(t, config.Port, 486)
		os.Unsetenv("EMAIL_USERNAME")
		os.Unsetenv("EMAIL_PASSWORD")
		os.Unsetenv("EMAIL_HOST")
		os.Unsetenv("EMAIL_PORT")
	})

	t.Run("MissingOptionalEnvVarsReturnsDefaults", func(t *testing.T) {
		os.Setenv("EMAIL_USERNAME", "some@email.io")
		os.Setenv("EMAIL_PASSWORD", "abc1234")
		config, err := getEmailConfig()
		expect.NoError(t, err)
		expect.Equal(t, config.Host, DefaultHost)
		expect.Equal(t, config.Port, DefaultPort)
		os.Unsetenv("EMAIL_USERNAME")
		os.Unsetenv("EMAIL_PASSWORD")
	})
}

func TestCreateAndLogReport(t *testing.T) {
	report := &Report{}

	t.Run("AddSuccessfulCheckIncreasesSuccessfulCount", func(t *testing.T) {
		report.Add(&sermoncore.ServiceStatus{
			Name:    "good.test",
			Healthy: true,
			Err:     nil,
		})
		report.Add(&sermoncore.ServiceStatus{
			Name:    "alsogood.test",
			Healthy: true,
			Err:     nil,
		})
		expect.Equal(t, report.Successful, 2)
	})

	t.Run("AddFailedCheckIncreasesFailedCount", func(t *testing.T) {
		report.Add(&sermoncore.ServiceStatus{
			Name:    "bad.test",
			Healthy: false,
			Err:     errors.New(""),
		})
		expect.Equal(t, report.Failed, 1)
	})

	t.Run("PrintsReportToGivenWriter", func(t *testing.T) {
		var buf bytes.Buffer
		report.Log(&buf)
		reportStr := buf.String()

		expect.Contains(t, reportStr, "SUCCESSFUL: 2")
		expect.Contains(t, reportStr, "FAILED: 1")
		expect.Contains(t, reportStr, "TOTAL: 3")
		expect.Contains(t, reportStr, "GET good.test -> OK")
		expect.Contains(t, reportStr, "GET bad.test -> ERROR")
	})
}

func TestSome(t *testing.T) {
	report := &Report{
		Services: []*sermoncore.ServiceStatus{
			{Name: "one", Healthy: true, Err: nil},
			{Name: "two", Healthy: true, Err: nil},
		},
	}

	t.Run("ReturnsFalseWhenNoElementSatisfiesCondition", func(t *testing.T) {
		hasUnhealthy := some(report.Services, func(ss *sermoncore.ServiceStatus) bool {
			return !ss.Healthy
		})
		expect.Equal(t, hasUnhealthy, false)
	})

	t.Run("ReturnsTrueWhenAllElementsSatisfyCondition", func(t *testing.T) {
		hasHealthy := some(report.Services, func(ss *sermoncore.ServiceStatus) bool {
			return ss.Healthy
		})
		expect.Equal(t, hasHealthy, true)
	})

	t.Run("ReturnsTrueWhenOneElementSatisfiesCondition", func(t *testing.T) {
		hasOne := some(report.Services, func(ss *sermoncore.ServiceStatus) bool {
			return ss.Name == "one"
		})
		expect.Equal(t, hasOne, true)
	})
}
