package xlog

import (
	"testing"
	"reflect"
	"strings"
	"fmt"
	"io"
)

// LoggerName is the default name for the test logger.
const LoggerName = "testing"

// Fixture creates and returns a new logger and writer.
func LoggerFixture(level Level) (*DefaultLogger, *MemoryWriter) {
	writer := NewMemoryWriter()
	logger := NewLogger(LoggerName)
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

// ActualIsNotEmpty assets that actual is not an empty string.
func ActualIsNotEmpty(t *testing.T, actual string) {
	if actual != "" {
		t.Errorf("Expected empty string but got '%s'.", actual)
	}
}

// TestName -
func TestName(t *testing.T) {
	logger, _ := LoggerFixture(DebugLevel)
	if logger.Name() != LoggerName {
		t.Errorf("Expected a logger named '%s'.", LoggerName)
	}
}

// TestLevelCalls -
func TestLevelCalls(t *testing.T) {
	logger, writer := LoggerFixture(DebugLevel)
	var methods = map[Level]string{
		DebugLevel: "Debug",
		InfoLevel: "Info",
		NoticeLevel: "Notice",
		WarningLevel: "Warning",
		ErrorLevel: "Error",
		CriticalLevel: "Critical",
		AlertLevel: "Alert",
		EmergencyLevel: "Emergency",
	}

	for level, method := range methods {
		writer.Clear()
		Invoke(logger, method, "This is a test.")
		expected := fmt.Sprintf("testing.%s This is a test.", Levels[level])
		ActualContains(t, writer.String(), expected)
	}
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
	ActualIsNotEmpty(t, writer.String())
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
	ActualIsNotEmpty(t, writer.String())

	writer.Clear()
	logger.Enabled = true
	logger.Debug(expected)
	ActualContains(t, writer.String(), expected)
}

// TestAliases -
func TestAliases(t *testing.T) {
	writer := NewMemoryWriter()
	Aliases["stdout"] = writer

	logger := NewLogger(LoggerName)
	logger.Append("stdout", DebugLevel)
	logger.Debug("This is a test.")
	expected := "testing.DEBUG This is a test."
	ActualContains(t, writer.String(), expected)
}

// TestDefaults -
func TestDefaults(t *testing.T) {
	writer := NewMemoryWriter()
	DefaultMessageFormat = "{level} - {message}"
	DefaultAppendLevel = WarningLevel
	DefaultAppendWriters = []io.Writer{writer}
	logger := NewLogger(LoggerName)

	logger.Warning("This is a test.")
	expected := "WARNING - This is a test.\n"
	actual := writer.String()
	if expected != actual {
		t.Errorf("Expected '%s' to equal '%s'.", expected, actual)
	}

	writer.Clear()
	logger.Debug("This is a test.")
	ActualIsNotEmpty(t, writer.String())
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

// Invoke calls the named method on any interface with the given arguments.
func Invoke(any interface{}, name string, args... interface{}) {
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
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
