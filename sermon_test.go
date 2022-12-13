package sermon

import (
	"testing"

	"gitlab.com/germandv/sermon/expect"
)

func TestWithRetry(t *testing.T) {
	t.Run("FailingFunctionGetsCalledNTimes", func(t *testing.T) {
		t.Parallel()
		maxAttempts := 5
		wantAttempts := 5
		invocations := 0

		fn := withRetry(
			maxAttempts,
			func(n int) *int {
				invocations++ // keep track of how many times the function is called.
				return &n     // return value is not important for the test.
			},
			func(n *int) bool {
				return true // `true` means we should retry.
			})

		fn(9) // the number is not important for the test.
		expect.Equal(t, invocations, wantAttempts)
	})

	t.Run("PassingFunctionGetsCalledOneTimeOnly", func(t *testing.T) {
		t.Parallel()
		maxAttempts := 5
		wantAttempts := 1
		invocations := 0

		fn := withRetry(
			maxAttempts,
			func(n int) *int {
				invocations++ // keep track of how many times the function is called.
				return &n     // return value is not important for the test.
			},
			func(n *int) bool {
				return false // `false` means we should not re-run the function.
			})

		fn(9) // the number is not important for the test.
		expect.Equal(t, invocations, wantAttempts)
	})
}
