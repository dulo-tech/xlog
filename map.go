package xlog

import "log"

// LoggerMap is an interface that stores a map of log levels and loggers.
type LoggerMap interface {
	Append(logger *log.Logger, level Level)
	FindByLevel(level Level) []*log.Logger
}

// Loggers maps loggers to levels.
type DefaultLoggerMap struct {
	loggers map[Level][]*log.Logger
}

// NewDefaultLoggerMap creates and returns a *DefaultLoggerMap instance.
func NewDefaultLoggerMap() *DefaultLoggerMap {
	loggers := make(map[Level][]*log.Logger, len(Levels))
	for level, _ := range Levels {
		loggers[level] = make([]*log.Logger, 0, InitialLoggerCapacity)
	}

	return &DefaultLoggerMap{loggers}
}

// Append adds a logger to the map at the given level.
func (m *DefaultLoggerMap) Append(logger *log.Logger, level Level) {
	for lev, _ := range m.loggers {
		if lev >= level {
			m.loggers[lev] = append(m.loggers[lev], logger)
		}
	}
}

// FindLevel returns the loggers at the given level or higher.
func (m *DefaultLoggerMap) FindByLevel(level Level) []*log.Logger {
	return m.loggers[level]
}
