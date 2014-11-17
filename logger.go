package xlog

import (
	"log"
	"os"
	"fmt"
	"io"
)

// DefaultLogger is the default implementation of the Loggable interface.
type DefaultLogger struct {
	// Enabled defines whether logging is enabled.
	Enabled  bool

	// Formatter is used to format the log messages.
	Formatter Formatter

	// Loggers holds the appended file loggers.
	Loggers LoggerMap

	// name defines the name of the logger.
	name string

	// pointers contains any files that have been opened for logging.
	pointers []*os.File

	// closed defines whether the logger has been closed.
	closed bool
}

// NewLogger returns a *DefaultLogger instance that's been initialized with default values.
func NewLogger(name string) *DefaultLogger {
	logger := &DefaultLogger{
		name: name,
		Enabled: true,
		Formatter: NewDefaultFormatter(DefaultMessageFormat),
		Loggers: NewDefaultLoggerMap(),
		pointers: make([]*os.File, 0),
		closed: false,
	}
	if DefaultAppendFiles != nil && len(DefaultAppendFiles) > 0 {
		logger.MultiAppend(DefaultAppendFiles, DefaultAppendLevel)
	}
	if DefaultAppendWriters != nil && len(DefaultAppendWriters) > 0 {
		logger.MultiAppendWriters(DefaultAppendWriters, DefaultAppendLevel)
	}

	return logger
}

// NewMultiLogger returns a *DefaultLogger instance that's been initialized with one or
// more files at the given level.
func NewMultiLogger(name string, files []string, level Level) *DefaultLogger {
	logger := NewLogger(name)
	logger.MultiAppend(files, level);
	return logger;
}

// NewMultiWriterLogger returns a *DefaultLogger instance that's been initialized with one or
// more writers at the given level.
func NewMultiWriterLogger(name string, writers []io.Writer, level Level) *DefaultLogger {
	logger := NewLogger(name)
	logger.MultiAppendWriters(writers, level);
	return logger;
}

// NewFormattedLogger returns a *DefaultLogger instance using the provided formatter.
func NewFormattedLogger(formatter Formatter) *DefaultLogger {
	logger := NewLogger("")
	logger.Formatter = formatter
	return logger
}

// Name returns the name of the logger.
func (l *DefaultLogger) Name() string {
	return l.name
}

// Writable returns true when logging is enabled, and the logger hasn't been closed.
func (l *DefaultLogger) Writable() bool {
	return l.Enabled && !l.closed
}

// Closed returns whether the logger has been closed.
func (l *DefaultLogger) Closed() bool {
	return l.closed
}

// Close disables logging and frees up resources used by the logger.
// Note this method only closes files opened by the logger. It's the user's
// responsibility to close files that were passed to the logger via the
// AppendWriter method.
func (l *DefaultLogger) Close() {
	if !l.closed {
		for _, pointer := range l.pointers {
			pointer.Close()
		}

		l.Enabled = false
		l.Loggers = nil
		l.pointers = nil
		l.closed = true
	}
}

// Append adds a file that will be written to at the given level or greater.
// The file argument may be either the full path to a system file, or one of the
// aliases "stdout", "stdin", or "stderr".
func (l *DefaultLogger) Append(file string, level Level) {
	if w, ok := Aliases[file]; ok {
		l.Loggers.Append(newLogger(w), level)
	} else {
		w := l.open(file)
		if w != nil {
			l.Loggers.Append(newLogger(w), level)
			l.pointers = append(l.pointers, w)
		}
	}
}

// MultiAppend adds one or more files to the logger.
func (l *DefaultLogger) MultiAppend(files []string, level Level) {
	for _, file := range files {
		l.Append(file, level)
	}
}

// AppendWriter adds a writer that will be written to at the given level or greater.
func (l *DefaultLogger) AppendWriter(writer io.Writer, level Level) {
	l.Loggers.Append(newLogger(writer), level)
}

// MultiAppendWriters adds one or more io.Writer instances to the logger.
func (l *DefaultLogger) MultiAppendWriters(writers []io.Writer, level Level) {
	for _, writer := range writers {
		l.AppendWriter(writer, level)
	}
}

// ClearAppended removes all the files that have been appended to the logger.
func (l *DefaultLogger) ClearAppended() {
	l.Loggers.Clear()
}

// Log writes the message to each logger appended at the given level or higher.
// Arguments are handled in the manner of fmt.Print.
func (l *DefaultLogger) Log(level Level, v ...interface{}) {
	if l.Writable() {
		message := l.Formatter.Format(l.name, level, v...)
		if message != "" {
			for _, logger := range l.Loggers.FindByLevel(level) {
				logger.Print(message)
			}
		}

		if FatalOn&level > 0 {
			os.Exit(1)
		} else if PanicOn&level > 0 {
			panic(message)
		}
	}
}

// Logf writes the message to each logger appended at the given level or higher.
// Arguments are handled in the manner of fmt.Printf.
func (l *DefaultLogger) Logf(level Level, format string, v ...interface{}) {
	l.Log(level, fmt.Sprintf(format, v...))
}

// Debug writes to the logger at DebugLevel.
// Arguments are handled in the manner of fmt.Print.
func (l *DefaultLogger) Debug(v ...interface{}) {
	l.Log(DebugLevel, v...)
}

// Debugf writes to the logger at DebugLevel.
// Arguments are handled in the manner of fmt.Printf.
func (l *DefaultLogger) Debugf(format string, v ...interface{}) {
	l.Logf(DebugLevel, format, v...)
}

// Info writes to the logger at InfoLevel.
// Arguments are handled in the manner of fmt.Print.
func (l *DefaultLogger) Info(v ...interface{}) {
	l.Log(InfoLevel, v...)
}

// Infof writes to the logger at InfoLevel.
// Arguments are handled in the manner of fmt.Printf.
func (l *DefaultLogger) Infof(format string, v ...interface{}) {
	l.Logf(InfoLevel, format, v...)
}

// Notice writes to the logger at NoticeLevel.
// Arguments are handled in the manner of fmt.Print.
func (l *DefaultLogger) Notice(v ...interface{}) {
	l.Log(NoticeLevel, v...)
}

// Noticef writes to the logger at NoticeLevel.
// Arguments are handled in the manner of fmt.Printf.
func (l *DefaultLogger) Noticef(format string, v ...interface{}) {
	l.Logf(NoticeLevel, format, v...)
}

// Warning writes to the logger at WarningLevel.
// Arguments are handled in the manner of fmt.Print.
func (l *DefaultLogger) Warning(v ...interface{}) {
	l.Log(WarningLevel, v...)
}

// Warningf writes to the logger at WarningLevel.
// Arguments are handled in the manner of fmt.Printf.
func (l *DefaultLogger) Warningf(format string, v ...interface{}) {
	l.Logf(WarningLevel, format, v...)
}

// Error writes to the logger at ErrorLevel.
// Arguments are handled in the manner of fmt.Print.
func (l *DefaultLogger) Error(v ...interface{}) {
	l.Log(ErrorLevel, v...)
}

// Errorf writes to the logger at ErrorLevel.
// Arguments are handled in the manner of fmt.Printf.
func (l *DefaultLogger) Errorf(format string, v ...interface{}) {
	l.Logf(ErrorLevel, format, v...)
}

// Critical writes to the logger at CriticalLevel.
// Arguments are handled in the manner of fmt.Print.
func (l *DefaultLogger) Critical(v ...interface{}) {
	l.Log(CriticalLevel, v...)
}

// Criticalf writes to the logger at CriticalLevel.
// Arguments are handled in the manner of fmt.Printf.
func (l *DefaultLogger) Criticalf(format string, v ...interface{}) {
	l.Logf(CriticalLevel, format, v...)
}

// Alert writes to the logger at AlertLevel.
// Arguments are handled in the manner of fmt.Print.
func (l *DefaultLogger) Alert(v ...interface{}) {
	l.Log(AlertLevel, v...)
}

// Alertf writes to the logger at AlertLevel.
// Arguments are handled in the manner of fmt.Printf.
func (l *DefaultLogger) Alertf(format string, v ...interface{}) {
	l.Logf(AlertLevel, format, v...)
}

// Emergency writes to the logger at EmergencyLevel.
// Arguments are handled in the manner of fmt.Print.
func (l *DefaultLogger) Emergency(v ...interface{}) {
	l.Log(EmergencyLevel, v...)
}

// Emergencyf writes to the logger at EmergencyLevel.
// Arguments are handled in the manner of fmt.Printf.
func (l *DefaultLogger) Emergencyf(format string, v ...interface{}) {
	l.Logf(EmergencyLevel, format, v...)
}

// open returns a file that logs can be written to.
func (l *DefaultLogger) open(name string) *os.File {
	w, err := os.OpenFile(name, FileOpenFlags, FileOpenMode)
	if err != nil {
		if PanicOnFileErrors {
			panic(err)
		} else {
			w = nil
		}
	}

	return w
}

// newLogger returns a *log.Logger instance configured with the default options.
func newLogger(w io.Writer) *log.Logger {
	return log.New(w, "", 0)
}

