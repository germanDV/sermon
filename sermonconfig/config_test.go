package sermonconfig

import (
	"testing"

	"gitlab.com/germandv/sermon/expect"
)

func TestParse_BadEmail(t *testing.T) {
	t.Parallel()
	config, err := Parse(expect.ReadFile(t, "bad_email.toml"))
	expect.Nil(t, config)
	expect.Contains(t, err.Error(), "Invalid email address")
}

func TestParse_ZeroAttempts(t *testing.T) {
	t.Parallel()
	config, err := Parse(expect.ReadFile(t, "zero_attempts.toml"))
	expect.Nil(t, config)
	expect.Contains(t, err.Error(), "Invalid number of attempts")
}

func TestParse_BadURL(t *testing.T) {
	t.Parallel()
	config, err := Parse(expect.ReadFile(t, "bad_url.toml"))
	expect.Nil(t, config)
	expect.Contains(t, err.Error(), "invalid URI")
}

func TestParse_BadStatusCode(t *testing.T) {
	t.Parallel()
	config, err := Parse(expect.ReadFile(t, "bad_codes.toml"))
	expect.Nil(t, config)
	expect.Contains(t, err.Error(), "Invalid status code")
}

func TestParse_BadTimeout(t *testing.T) {
	t.Parallel()
	config, err := Parse(expect.ReadFile(t, "bad_timeout.toml"))
	expect.Nil(t, config)
	expect.Contains(t, err.Error(), "missing unit in duration")
}

func TestParse_MissingURL(t *testing.T) {
	t.Parallel()
	config, err := Parse(expect.ReadFile(t, "missing_url.toml"))
	expect.Nil(t, config)
	expect.Contains(t, err.Error(), "Missing `endpoint`")
}

func TestParse_MissingStatusCodes(t *testing.T) {
	t.Parallel()
	config, err := Parse(expect.ReadFile(t, "missing_codes.toml"))
	expect.Nil(t, config)
	expect.Contains(t, err.Error(), "Missing `codes`")
}

func TestParse_MissingTimeout(t *testing.T) {
	t.Parallel()
	config, err := Parse(expect.ReadFile(t, "missing_timeout.toml"))
	expect.Nil(t, config)
	expect.Contains(t, err.Error(), "Missing `timeout`")
}

func TestParse_MissingEmail(t *testing.T) {
	t.Parallel()
	config, err := Parse(expect.ReadFile(t, "missing_email.toml"))
	expect.Nil(t, config)
	expect.Contains(t, err.Error(), "Missing `email`")
}

func TestParse_MissingAttempts(t *testing.T) {
	t.Parallel()
	config, err := Parse(expect.ReadFile(t, "missing_attempts.toml"))
	expect.Nil(t, config)
	expect.Contains(t, err.Error(), "Missing `attempts`")
}

func TestParse_GoodConfig(t *testing.T) {
	t.Parallel()
	_, err := Parse(expect.ReadFile(t, "good.toml"))
	expect.NoError(t, err)
}
