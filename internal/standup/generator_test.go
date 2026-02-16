package standup

import (
	"testing"
	"time"
)

func TestSinceDate(t *testing.T) {
	// SinceDate uses time.Now() so we test the logic conceptually
	// by verifying it returns a valid date string
	result := SinceDate()
	_, err := time.Parse("2006-01-02", result)
	if err != nil {
		t.Fatalf("SinceDate returned invalid date %q: %v", result, err)
	}

	now := time.Now()
	parsed, _ := time.Parse("2006-01-02", result)

	if now.Weekday() == time.Monday {
		// Should be 3 days ago (Friday)
		expected := now.AddDate(0, 0, -3).Format("2006-01-02")
		if result != expected {
			t.Errorf("on Monday, expected %s, got %s", expected, result)
		}
	} else {
		// Should be 1 day ago
		expected := now.AddDate(0, 0, -1).Format("2006-01-02")
		if result != expected {
			t.Errorf("expected %s, got %s", expected, result)
		}
	}

	// The date should be in the past
	if !parsed.Before(now) {
		t.Errorf("SinceDate should return a past date, got %s", result)
	}
}
