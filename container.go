package xlog

import "log"

// Container is an interface that stores a container of log levels and loggers.
type Container interface {
	Append(logger *log.Logger, level Level)
	Get(level Level) []*log.Logger
	Clear()
}

// DefaultContainer maps loggers to levels.
type DefaultContainer struct {
	Capacity int
	loggers map[Level][]*log.Logger
}

// NewDefaultContainer creates and returns a *DefaultLoggerContainer instance.
func NewDefaultContainer(capacity int) *DefaultContainer {
	lm := &DefaultContainer{Capacity: capacity, loggers: nil}
	lm.Clear()
	return lm
}

// Append adds a logger to the container at the given level.
func (m *DefaultContainer) Append(logger *log.Logger, level Level) {
	for lev, _ := range m.loggers {
		if (lev&level > 0) || (lev >= level) {
			m.loggers[lev] = append(m.loggers[lev], logger)
		}
	}
}

// Get returns the loggers at the given level or higher.
func (m *DefaultContainer) Get(level Level) []*log.Logger {
	return m.loggers[level]
}

// Clear removes all the appended loggers.
func (m *DefaultContainer) Clear() {
	m.loggers = make(map[Level][]*log.Logger, len(Levels))
	for level, _ := range Levels {
		m.loggers[level] = make([]*log.Logger, 0, m.Capacity)
	}
}
