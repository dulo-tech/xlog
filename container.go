package xlog

import (
	"io"
	"log"
	"os"
)

// Container is an interface that stores a container of log levels and loggers.
type Container interface {
	Append(writer io.Writer, level Level)
	Get(level Level) []*log.Logger
	Clear()
	Close()
	Closed() bool
}

// DefaultContainer maps loggers to levels.
type DefaultContainer struct {
	// Capacity is the initial number of loggers to make.
	Capacity int

	// loggers are the loggers to be written to.
	loggers map[Level][]*log.Logger

	// pointers contains any files that have been opened for logging.
	pointers []*os.File

	// closed defines whether the logger has been closed.
	closed bool
}

// NewDefaultContainer creates and returns a *DefaultLoggerContainer instance.
func NewDefaultContainer(capacity int) *DefaultContainer {
	lm := &DefaultContainer{
		Capacity: capacity,
		loggers:  nil,
		pointers: make([]*os.File, 0, DefaultInitialCapacity),
		closed:   false,
	}
	lm.Clear()

	return lm
}

// Append adds a logger to the container at the given level.
func (m *DefaultContainer) Append(writer io.Writer, level Level) {
	logger := newLogger(writer)
	for lev := range m.loggers {
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
	for level := range Levels {
		m.loggers[level] = make([]*log.Logger, 0, m.Capacity)
	}
}

// Close closes any resources being used by the container.
func (m *DefaultContainer) Close() {
	if !m.closed {
		for _, pointer := range m.pointers {
			pointer.Close()
		}
		m.pointers = nil
		m.closed = true
	}
}

// Closed returns whether the container has been closed.
func (m *DefaultContainer) Closed() bool {
	return m.closed
}

// newLogger returns a *log.Logger instance configured with the default options.
func newLogger(writer io.Writer) *log.Logger {
	return log.New(writer, "", 0)
}
