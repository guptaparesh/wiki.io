package searcher

import (
	"testing"
)

func TestRunQuery(t *testing.T) {
	if actual := RunQuery("Seattle"); actual == nil {
		t.Errorf("RunQuery() Actual= %q, does not contain any search hits", actual)
	}
}

func TestFindMostViewed(t *testing.T) {
	const n = 5
	if actual := FindMostViewed(n); len(actual) != n {
		t.Errorf("FindMostViewed(%d) Actual number of results: %d, expected %d", n, len(actual), n)
	}
}