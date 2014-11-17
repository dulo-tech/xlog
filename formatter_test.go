package xlog

import "testing"

// TestFormat -
func TestFormat(t *testing.T) {
	formatter := NewDefaultFormatter("{level} {message}")
	actual := formatter.Format("testing", DebugLevel, "This is a test.")
	expected := "DEBUG This is a test."
	ActualEquals(t, actual, expected)
}

// TestPlaceholderFunc -
func TestPlaceholderFunc(t *testing.T) {
	formatter := NewDefaultFormatter("{level} {hostname} {message}")
	formatter.PlaceholderFunc("hostname", func(key string) string {
			return "test-service"
		})
	actual := formatter.Format("testing", DebugLevel, "This is a test.")
	expected := "DEBUG test-service This is a test."
	ActualEquals(t, actual, expected)
}
