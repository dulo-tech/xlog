package xlog

import (
	"os"
	"io"
)

// Level describes a logging level.
type Level int

const (
	DebugLevel Level = 1 << iota
	InfoLevel        = 1 << iota
	NoticeLevel      = 1 << iota
	WarningLevel     = 1 << iota
	ErrorLevel       = 1 << iota
	CriticalLevel    = 1 << iota
	AlertLevel       = 1 << iota
	EmergencyLevel   = 1 << iota
)

// levelOrder defines the order of the levels.
var levelOrder = []Level{
	DebugLevel,
	InfoLevel,
	NoticeLevel,
	WarningLevel,
	ErrorLevel,
	CriticalLevel,
	AlertLevel,
	EmergencyLevel,
}

// Levels maps Level to a string representation.
var Levels = map[Level]string{
	DebugLevel: "DEBUG",
	InfoLevel: "INFO",
	NoticeLevel: "NOTICE",
	WarningLevel: "WARNING",
	ErrorLevel: "ERROR",
	CriticalLevel: "CRITICAL",
	AlertLevel: "ALERT",
	EmergencyLevel: "EMERGENCY",
}

// IsGreaterLevel returns whether the level is_greater_than is greater than that.
func IsGreaterLevel(is_greater_than, that Level) bool {
	return searchForLevel(is_greater_than) > searchForLevel(that)
}

// IsLesserLevel returns whether the level is_less_than is less than that.
func IsLesserLevel(is_less_than, that Level) bool {
	return searchForLevel(is_less_than) < searchForLevel(that)
}

// searchForLevel returns the index for the given level or -1.
func searchForLevel(level Level) int {
	for idx, val := range levelOrder {
		if val == level {
			return idx
		}
	}
	panic("Invalid level.")
}

// Aliases maps file aliases to real file pointers.
var Aliases = map[string]io.Writer{
	"stdout": os.Stdout,
	"stdin": os.Stdin,
	"stderr": os.Stderr,
}

// Loggable is an interface that provides methods for logging messages to
// various levels.
type Loggable interface {
	Writable() bool
	Closed() bool
	Log(level Level, v ...interface{})
	Logf(level Level, format string, v ...interface{})
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Notice(v ...interface{})
	Noticef(format string, v ...interface{})
	Warning(v ...interface{})
	Warningf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Critical(v ...interface{})
	Criticalf(format string, v ...interface{})
	Alert(v ...interface{})
	Alertf(format string, v ...interface{})
	Emergency(v ...interface{})
	Emergencyf(format string, v ...interface{})
	Writer(level Level) *LoggerWriter
}
