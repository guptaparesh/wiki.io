package search

import "testing"

func TestHello(t *testing.T) {
	exp := "Hello, world."
	if actual := Hello(); actual != exp {
		t.Errorf("Hello() Actual= %q, expected = %q", actual, exp)
	}
}