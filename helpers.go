package sermon

import "net/http"

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
