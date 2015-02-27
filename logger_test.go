package xlog

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// LoggerName is the default name for the test logger.
const LoggerName = "testing"

// Fixture creates and returns a new logger and writer.
func LoggerFixture(level Level) (*DefaultLogger, *MemoryWriter) {
	writer := NewMemoryWriter()
	logger := New(LoggerName)
	logger.AppendWriter(writer, level)

	return logger, writer
}

// ActualEquals asserts that actual equals expected.
func ActualEquals(t *testing.T, actual, expected string) {
	if actual != expected {
		t.Errorf("Expected '%s' but got '%s'.", expected, actual)
	}
}

// ActualContains asserts that actual contains the sub-string expected.
func ActualContains(t *testing.T, actual, expected string) {
	if !strings.Contains(actual, expected) {
		t.Errorf("Expected actual '%s' to contain expected '%s'.", actual, expected)
	}
}

// ActualIsEmpty assets that actual is an empty string.
func ActualIsEmpty(t *testing.T, actual string) {
	if actual != "" {
		t.Errorf("Expected empty string but got '%s'.", actual)
	}
}

// ActualIsNotEmpty assets that actual is not an empty string
func ActualIsNotEmpty(t *testing.T, actual string) {
	if actual == "" {
		t.Error("Expected non-empty string")
	}
}

// ActualLevelGreaterThan asserts the actual level is greater than the expected.
func ActualLevelGreaterThan(t *testing.T, actual, expected Level) {
	if !IsGreaterLevel(actual, expected) {
		t.Errorf("Expected level %s to be greater than %s.", Levels[actual], Levels[expected])
	}
}

// ActualLevelLessThan asserts the actual level is less than the expected.
func ActualLevelLessThan(t *testing.T, actual, expected Level) {
	if !IsLesserLevel(actual, expected) {
		t.Errorf("Expected level %s to be less than %s.", Levels[actual], Levels[expected])
	}
}

// TestName -
func TestName(t *testing.T) {
	logger, _ := LoggerFixture(DebugLevel)
	if logger.Name != LoggerName {
		t.Errorf("Expected a logger named '%s'.", LoggerName)
	}
}

// TestGreater -
func TestGreater(t *testing.T) {
	ActualLevelGreaterThan(t, EmergencyLevel, DebugLevel)
	ActualLevelGreaterThan(t, AlertLevel, InfoLevel)
	ActualLevelGreaterThan(t, InfoLevel, DebugLevel)
	ActualLevelLessThan(t, DebugLevel, ErrorLevel)
	ActualLevelLessThan(t, InfoLevel, WarningLevel)
	ActualLevelLessThan(t, AlertLevel, EmergencyLevel)
}

type LevelCall struct {
	method  string
	greater []Level
	lower   []Level
}

// TestLevelCalls -
func TestLevelCalls(t *testing.T) {
	var calls = map[Level]LevelCall{
		DebugLevel: {
			"Debug",
			[]Level{InfoLevel, NoticeLevel, WarningLevel, ErrorLevel, CriticalLevel, AlertLevel, EmergencyLevel},
			[]Level{},
		},
		InfoLevel: {
			"Info",
			[]Level{NoticeLevel, WarningLevel, ErrorLevel, CriticalLevel, AlertLevel, EmergencyLevel},
			[]Level{DebugLevel},
		},
		NoticeLevel: {
			"Notice",
			[]Level{WarningLevel, ErrorLevel, CriticalLevel, AlertLevel, EmergencyLevel},
			[]Level{DebugLevel, InfoLevel},
		},
		WarningLevel: {
			"Warning",
			[]Level{ErrorLevel, CriticalLevel, AlertLevel, EmergencyLevel},
			[]Level{DebugLevel, InfoLevel, NoticeLevel},
		},
		ErrorLevel: {
			"Error",
			[]Level{CriticalLevel, AlertLevel, EmergencyLevel},
			[]Level{DebugLevel, InfoLevel, NoticeLevel, WarningLevel},
		},
		CriticalLevel: {
			"Critical",
			[]Level{AlertLevel, EmergencyLevel},
			[]Level{DebugLevel, InfoLevel, NoticeLevel, WarningLevel, ErrorLevel},
		},
		AlertLevel: {
			"Alert",
			[]Level{EmergencyLevel},
			[]Level{DebugLevel, InfoLevel, NoticeLevel, WarningLevel, ErrorLevel, CriticalLevel},
		},
		EmergencyLevel: {
			"Emergency",
			[]Level{},
			[]Level{DebugLevel, InfoLevel, NoticeLevel, WarningLevel, ErrorLevel, CriticalLevel, AlertLevel},
		},
	}
	for level, call := range calls {
		logger, writer := LoggerFixture(level)

		Invoke(logger, call.method, "This is a test.")
		expected := fmt.Sprintf("testing.%s This is a test.", Levels[level])
		ActualContains(t, writer.String(), expected)

		for _, greater_level := range call.greater {
			writer.Clear()
			Invoke(logger, calls[greater_level].method, "This is a test.")
			expected := fmt.Sprintf("testing.%s This is a test.", Levels[greater_level])
			ActualContains(t, writer.String(), expected)
		}
		for _, lower_level := range call.lower {
			writer.Clear()
			Invoke(logger, calls[lower_level].method, "This is a test.")
			ActualIsEmpty(t, writer.String())
		}
	}

	logger, writer := LoggerFixture(WarningLevel | InfoLevel)
	logger.Debug("This is a test.")
	ActualIsEmpty(t, writer.String())

	writer.Clear()
	logger.Info("This is a test.")
	ActualIsNotEmpty(t, writer.String())

	writer.Clear()
	logger.Notice("This is a test.")
	ActualIsEmpty(t, writer.String())

	writer.Clear()
	logger.Warning("This is a test.")
	ActualIsNotEmpty(t, writer.String())

	writer.Clear()
	logger.Error("This is a test.")
	ActualIsNotEmpty(t, writer.String())
}

// TestLevels -
func TestLevels(t *testing.T) {
	logger, writer := LoggerFixture(WarningLevel)

	logger.Warning("This is a test.")
	expected := "testing.WARNING This is a test."
	ActualContains(t, writer.String(), expected)

	logger.Critical("This is a test.")
	expected = "testing.CRITICAL This is a test."
	ActualContains(t, writer.String(), expected)

	writer.Clear()
	logger.Debug("This is a test.")
	ActualIsEmpty(t, writer.String())
}

// TestWriter -
func TestWriter(t *testing.T) {
	logger, writer := LoggerFixture(DebugLevel)
	lw := logger.Writer(DebugLevel)

	fmt.Fprint(lw, "This is a test.")
	expected := "testing.DEBUG This is a test."
	ActualContains(t, writer.String(), expected)
}

// TestEnabled -
func TestEnabled(t *testing.T) {
	logger, writer := LoggerFixture(DebugLevel)
	expected := "testing.DEBUG This is a test."

	logger.Debug("This is a test.")
	ActualContains(t, writer.String(), expected)

	writer.Clear()
	logger.Enabled = false
	logger.Debug("This is a test.")
	ActualIsEmpty(t, writer.String())

	writer.Clear()
	logger.Enabled = true
	logger.Debug(expected)
	ActualContains(t, writer.String(), expected)
}

// TestAliases -
func TestAliases(t *testing.T) {
	writer := NewMemoryWriter()
	Aliases["stdout"] = writer

	logger := New(LoggerName)
	logger.Append("stdout", DebugLevel)
	logger.Debug("This is a test.")
	expected := "testing.DEBUG This is a test."
	ActualContains(t, writer.String(), expected)
}

// TestClose -
func TestClose(t *testing.T) {
	logger, _ := LoggerFixture(DebugLevel)
	if !logger.Writable() {
		t.Error("Expected Writable() to be true.")
	}

	logger.Close()
	if logger.Writable() {
		t.Error("Expected Writable() to be false.")
	}

	if !logger.Closed() {
		t.Error("Expected Closed() to return true.")
	}
}

// TestInstance -
func TestInstance(t *testing.T) {
	if Instance() != Instance() {
		t.Error("Expected the same instance from Instance().")
	}

	writer := NewMemoryWriter()
	AppendWriter(writer, DebugLevel)

	expected := "testing.DEBUG This is a test."
	Debug(expected)
	ActualContains(t, writer.String(), expected)
}

// TestGetLogger -
func TestGetLogger(t *testing.T) {
	loggerA := GetLogger("a")
	loggerB := GetLogger("b")

	if loggerA != GetLogger("a") {
		t.Error("Expected the same instance for loggerA.")
	}
	if loggerB != GetLogger("b") {
		t.Error("Expected the same instance for loggerB.")
	}
	if GetLogger("c") == GetLogger("d") {
		t.Error("Expected different instances for loggerC and loggerD.")
	}
}

// TestNewFromLogger -
func TestNewFromLogger(t *testing.T) {
	loggerA := GetLogger("a")
	loggerB := loggerA.New("b").(*DefaultLogger)
	if loggerA.Settings != loggerB.Settings {
		t.Error("Expected both loggers to have the same settings.")
	}
}

// Invoke calls the named method on any interface with the given arguments.
func Invoke(any interface{}, name string, args ...interface{}) {
	inputs := make([]reflect.Value, len(args))
	for i := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	reflect.ValueOf(any).MethodByName(name).Call(inputs)
}

// MemoryWriter -

type MemoryWriter struct {
	Data []byte
	Size int
}

func NewMemoryWriter() *MemoryWriter {
	return &MemoryWriter{}
}

func (w *MemoryWriter) Write(p []byte) (n int, err error) {
	w.Data = p
	w.Size = len(p)
	return w.Size, nil
}
func (w *MemoryWriter) String() string {
	return string(w.Data[:w.Size])
}
func (w *MemoryWriter) Clear() {
	w.Data = nil
	w.Size = 0
}
