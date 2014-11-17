package xlog

import "log"

// LoggerContainer is an interface that stores a container of log levels and loggers.
type LoggerContainer interface {
	Append(logger *log.Logger, level Level)
	FindByLevel(level Level) []*log.Logger
	Clear()
}

// DefaultLoggerContainer maps loggers to levels.
type DefaultLoggerContainer struct {
	loggers map[Level][]*log.Logger
}

// NewDefaultLoggerContainer creates and returns a *DefaultLoggerContainer instance.
func NewDefaultLoggerContainer() *DefaultLoggerContainer {
	lm := &DefaultLoggerContainer{}
	lm.Clear()
	return lm
}

// Append adds a logger to the container at the given level.
func (m *DefaultLoggerContainer) Append(logger *log.Logger, level Level) {
	for lev, _ := range m.loggers {
		if lev >= level {
			m.loggers[lev] = append(m.loggers[lev], logger)
		}
	}
}

// FindByLevel returns the loggers at the given level or higher.
func (m *DefaultLoggerContainer) FindByLevel(level Level) []*log.Logger {
	return m.loggers[level]
}

// Clear removes all the appended loggers.
func (m *DefaultLoggerContainer) Clear() {
	m.loggers = make(map[Level][]*log.Logger, len(Levels))
	for level, _ := range Levels {
		m.loggers[level] = make([]*log.Logger, 0, InitialLoggerCapacity)
	}
}
