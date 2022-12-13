package expect

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ReadFile reads a file from testdata/dir and returns the entire content as a string.
func ReadFile(t *testing.T, filename string) string {
	t.Helper()
	bytes, err := os.ReadFile(filepath.Join("..", "testdata", filename))
	if err != nil {
		t.Fatal(err)
	}
	return string(bytes)
}

// Nil fails the test if the element is not nil.
func Nil[T comparable](t *testing.T, element T) {
	t.Helper()
	zeroValue := new(T)
	if element != *zeroValue {
		t.Errorf("want nil, got %+v", element)
	}
}

// NoError fails the test if err != nil.
func NoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("want no error, got %s", err)
	}
}

// Contains fails the test if a given string does not contain a substring.
func Contains(t *testing.T, str string, sub string) {
	t.Helper()
	if !strings.Contains(str, sub) {
		t.Errorf("want %q to contain %q", str, sub)
	}
}

// Equal fails the test if a and b are not equal.
func Equal[T comparable](t *testing.T, a T, b T) {
	t.Helper()
	if a != b {
		t.Errorf("want %v to equal %v", a, b)
	}
}
