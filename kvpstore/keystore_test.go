package kvpstore

import (
	"log"
	"testing"
)

func TestExample(t *testing.T) {
	actual := Example("hallo", "Munchen")
	exp := "Munchen"
	if actual != exp {
		t.Errorf("Expected %s, actual %s", exp, actual)
	} else {
		log.Printf("Expected %s, actual %s\n", exp, actual)
	}
}