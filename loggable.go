package xlog

// Level describes a logging level.
type Level int

const (
	// Info useful to developers for debugging the application, not useful during operations.
	DebugLevel Level = 1 << iota
	
	// Normal operational messages - may be harvested for reporting, measuring
	// throughput, etc. - no action required.
	InfoLevel        = 1 << iota
	
	// Events that are unusual but not error conditions - might be summarized in an email to
	// developers or admins to spot potential problems - no immediate action required.
	NoticeLevel      = 1 << iota
	
	// Warning messages, not an error, but indication that an error will occur if action is not
	// taken, e.g. file system 85% full - each item must be resolved within a given time.
	WarningLevel     = 1 << iota
	
	// Non-urgent failures, these should be relayed to developers or admins; each item must be
	// resolved within a given time.
	ErrorLevel       = 1 << iota
	
	// Should be corrected immediately, but indicates failure in a secondary system, an example
	// is a loss of a backup ISP connection.
	CriticalLevel    = 1 << iota

	// Should be corrected immediately, therefore notify staff who can fix the problem. An
	// example would be the loss of a primary ISP connection.
	AlertLevel       = 1 << iota
	
	// A "panic" condition usually affecting multiple apps/servers/sites. At this level it
	// would usually notify all tech staff on call.
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
