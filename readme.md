# SERvice MONitor

Make requests to _health check_ endpoints and get notified when your services are down.

## Configuration

Configuration is set in `cmd/services.toml` file, please refer to `cmd/services.sample.toml` for an example.

This file is embeded into the binary, so you will need to update the code if you name it something different or move it to another location.

## Email

When at least one of the services is not healthy, an email will be sent to the `email` address specified in `services.toml`.

For this to work, you'll need to provide email server information to send the email from. This is done via environment variables:

- `EMAIL_USERNAME`: the _from_ email address. _Required_.
- `EMAIL_PASSWORD`: the password for the _from_ email address. _Required_.
- `EMAIL_HOST`: the email server host, defaults to `smtp.gmail.com`.
- `EMAIL_PORT`: the SMTP email server port, defaults to `587`.

## Usage

1. Rename `cmd/services.sample.toml` to `cmd/services.toml` and edit it with the service you wish to monitor.
1. Build a binary (ie: `go build -o bin/sermon cmd/main.go`)
1. Set up the required env vars for the email server.
1. Run the binary as a cron job with the desired frequency.

