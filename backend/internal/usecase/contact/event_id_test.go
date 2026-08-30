package contact

import (
	"regexp"
	"testing"
)

func TestNewEventIDReturnsUniqueUUIDv4(t *testing.T) {
	first, err := newEventID()
	if err != nil {
		t.Fatalf("newEventID() error = %v", err)
	}
	second, err := newEventID()
	if err != nil {
		t.Fatalf("newEventID() error = %v", err)
	}
	if first == second {
		t.Fatalf("newEventID() returned duplicate %q", first)
	}

	pattern := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	if !pattern.MatchString(first) {
		t.Fatalf("newEventID() = %q, want UUIDv4", first)
	}
}
