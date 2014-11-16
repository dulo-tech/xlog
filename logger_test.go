package xlog

import (
	"testing"
	"reflect"
	"strings"
	"fmt"
	"io"
)

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

func TestLevelCalls(t *testing.T) {
	writer := NewMemoryWriter()
	logger := NewLogger("testing")
	logger.AppendWriter(writer, DebugLevel)
	
	for level, method := range methods {
		writer.Clear()
		Invoke(logger, method, "This is a test.")
		
		expected := fmt.Sprintf("testing.%s This is a test.", Levels[level])
		actual := writer.String()
		if !strings.Contains(actual, expected) {
			t.Errorf("Expected '%s' to contain '%s'.", actual, expected)
		}
	}
}

func TestLevels(t *testing.T) {
	writer := NewMemoryWriter()
	logger := NewLogger("testing")
	logger.AppendWriter(writer, WarningLevel)

	expected := "testing.WARNING This is a test."
	logger.Warning(expected)
	actual := writer.String()
	if !strings.Contains(actual, expected) {
		t.Errorf("Expected '%s' to contain '%s'.", expected, actual)
	}

	expected = "testing.CRITICAL This is a test."
	logger.Critical(expected)
	actual = writer.String()
	if !strings.Contains(actual, expected) {
		t.Errorf("Expected '%s' to contain '%s'.", expected, actual)
	}
	
	writer.Clear()
	logger.Debug("This is a test.")
	actual = writer.String()
	if actual != "" {
		t.Errorf("Expected an empty string but got '%s'.", actual)
	}
}

func TestEnabled(t *testing.T) {
	writer := NewMemoryWriter()
	logger := NewLogger("testing")
	logger.AppendWriter(writer, DebugLevel)

	expected := "testing.DEBUG This is a test."
	logger.Debug(expected)
	actual := writer.String()
	if !strings.Contains(actual, expected) {
		t.Errorf("Expected '%s' to contain '%s'.", expected, actual)
	}

	writer.Clear()
	logger.Enabled = false
	logger.Debug("This is a test.")
	actual = writer.String()
	if actual != "" {
		t.Errorf("Expected an empty string but got '%s'.", actual)
	}

	writer.Clear()
	logger.Enabled = true
	logger.Debug(expected)
	actual = writer.String()
	if !strings.Contains(actual, expected) {
		t.Errorf("Expected '%s' to contain '%s'.", expected, actual)
	}
}

func TestAliases(t *testing.T) {
	writer := NewMemoryWriter()
	Aliases["stdout"] = writer
	
	logger := NewLogger("testing")
	logger.Append("stdout", DebugLevel)

	expected := "testing.DEBUG This is a test."
	logger.Debug(expected)
	actual := writer.String()
	if !strings.Contains(actual, expected) {
		t.Errorf("Expected '%s' to contain '%s'.", expected, actual)
	}
}

func TestDefaults(t *testing.T) {
	writer := NewMemoryWriter()
	DefaultMessageFormat = "{level} - {message}"
	DefaultAppendLevel = WarningLevel
	DefaultAppendWriters = []io.Writer{writer}
	logger := NewLogger("testing")

	logger.Warning("This is a test.")
	expected := "WARNING - This is a test.\n"
	actual := writer.String()
	if expected != actual {
		t.Errorf("Expected '%s' to equal '%s'.", expected, actual)
	}
	
	writer.Clear()
	logger.Debug("This is a test.")
	actual = writer.String()
	if actual != "" {
		t.Errorf("Expected an empty string but go '%s'.", actual)
	}
}

func TestClose(t *testing.T) {
	logger := NewLogger("testing")
	if !logger.Writable() {
		t.Error("Expected Writable() to be true.")
	}
	
	logger.Close()
	if logger.Writable() {
		t.Error("Expected Writable() to be false.")
	}
}

func TestInstance(t *testing.T) {
	if Instance() != Instance() {
		t.Error("Expected the same instance from Instance().")
	}
	
	writer := NewMemoryWriter()
	AppendWriter(writer, DebugLevel)

	expected := "testing.DEBUG This is a test."
	Debug(expected)
	actual := writer.String()
	if !strings.Contains(actual, expected) {
		t.Errorf("Expected '%s' to contain '%s'.", expected, actual)
	}
}

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


func Invoke(any interface{}, name string, args... interface{}) {
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	reflect.ValueOf(any).MethodByName(name).Call(inputs)
}

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
