package secret

import "fmt"

type Secret[T any] struct {
	value T
}

// New wraps the provided value in a `Secret` and returns it.
func New[T any](v T) Secret[T] {
	return Secret[T]{value: v}
}

func (s Secret[T]) Format(f fmt.State, verb rune) {
	fmt.Fprint(f, "Secret value access denied, call `Expose()` to read it.")
}

// Expose returns the wrapped secret value.
func (s Secret[T]) Expose() T {
	return s.value
}
