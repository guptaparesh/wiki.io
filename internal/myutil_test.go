package myutil

import (
	"strings"
	"testing"
)

func TestTrackTime(t *testing.T) {
	msg, _ := TrackTime("testing")
	if !strings.HasPrefix("testing", msg) {
		t.Errorf("expected to find msg with prefix 'testing', got %v", msg)
	}
}
