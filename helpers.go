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
