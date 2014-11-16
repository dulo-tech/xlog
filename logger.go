package xlog

import (
	"log"
	"os"
	"fmt"
	"io"
)

// DefaultDateFormat is the date format to use when none has been specified.
const DefaultDateFormat = "2006-01-02 15:04:05.000"

// DefaultMessageFormat is the message format to use when none has been specified.
const DefaultMessageFormat = "{date|" + DefaultDateFormat + "} {name}.{level} {message}"

// Level describes a logging level.
type Level int

const (
	Debug Level = 1 << iota
	Info        = 1 << iota
	Notice      = 1 << iota
	Warning     = 1 << iota
	Error       = 1 << iota
	Critical    = 1 << iota
	Alert       = 1 << iota
	Emergency   = 1 << iota
)

// Levels maps Level to a string representation.
var Levels = map[Level]string{
	Debug: "DEBUG",
	Info: "INFO",
	Notice: "NOTICE",
	Warning: "WARNING",
	Error: "ERROR",
	Critical: "CRITICAL",
	Alert: "ALERT",
	Emergency: "EMERGENCY",
}

// FileAliases maps file aliases to real file pointers.
var FileAliases = map[string]*os.File{
	"stdout": os.Stdout,
	"stdin": os.Stdin,
	"stderr": os.Stderr,
}

var (
	// FileFlags defines the file open options.
	FileFlags int = os.O_RDWR|os.O_CREATE | os.O_APPEND

	// FileMode defines the mode files are opened in.
	FileMode os.FileMode = 0666

	// PanicOnFileErrors defines whether the logger should panic when opening a file
	// fails. When set to false, any file open errors are ignored, and the file won't be
	// appended.
	PanicOnFileErrors = true

	// LoggerCapacity defines the initial capacity for each type of logger.
	LoggerCapacity = 2
)

// Loggable is an interface that provides methods for logging messages to
// various levels.
type Loggable interface {
	Log(level Level, v ...interface{})
	Logf(level Level, format string, v ...interface{})
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
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
}

// Logger is a light weight logger designed to write to multiple files at different
// log levels.
type Logger struct {
	// Enabled defines whether logging is enabled.
	Enabled  bool

	// Formatter is used to format the log messages.
	Formatter Formatter

	// Loggers holds the appended file loggers.
	Loggers LoggerMap

	// FatalOn represents levels that causes the application to exit.
	FatalOn Level

	// PanicOn represents levels that causes the application to panic.
	PanicOn Level

	// pointers contains any files that have been opened for logging.
	pointers []*os.File

	// closed defines whether the logger has been closed.
	closed bool
}

// NewLogger returns a *Logger instance that's been initialized with default values.
func NewLogger(name string) *Logger {
	return &Logger{
		Enabled: true,
		Formatter: NewDefaultFormatter(DefaultMessageFormat, name),
		Loggers: NewDefaultLoggerMap(),
		FatalOn: 0,
		PanicOn: 0,
		pointers: make([]*os.File, 0),
		closed: false,
	}
}

// NewFormattedLogger returns a *Logger instance using the provided formatter.
func NewFormattedLogger(formatter Formatter) *Logger {
	return &Logger{
		Enabled: true,
		Formatter: formatter,
		Loggers: NewDefaultLoggerMap(),
		FatalOn: 0,
		PanicOn: 0,
		pointers: make([]*os.File, 0),
		closed: false,
	}
}

// Append adds a file that will be written to at the given level or greater.
// The file argument may be either the full path to a system file, or one of the
// aliases "stdout", "stdin", or "stderr".
func (l *Logger) Append(file string, level Level) {
	if w, ok := FileAliases[file]; ok {
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
func (l *Logger) MultiAppend(files []string, level Level) {
	for _, file := range files {
		l.Append(file, level)
	}
}

// AppendWriter adds a writer that will be written to at the given level or greater.
func (l *Logger) AppendWriter(w io.Writer, level Level) {
	l.Loggers.Append(newLogger(w), level)
}

// MultiAppendWriter adds one or more io.Writer instances to the logger.
func (l *Logger) MultiAppendWriter(writers []io.Writer, level Level) {
	for _, writer := range writers {
		l.AppendWriter(writer, level)
	}
}

// Close disables logging and frees up resources used by the logger.
// Note this method only closes files opened by the logger. It's the user's
// responsibility to close files that were passed to the logger via the
// AppendWriter method.
func (l *Logger) Close() {
	if !l.closed {
		for _, pointer := range l.pointers {
			pointer.Close()
		}

		l.Enabled = false
		l.Loggers = nil
		l.pointers = nil
	}
}

// Writable returns true when logging is enabled, and the logger hasn't been closed.
func (l *Logger) Writable() bool {
	return l.Enabled && !l.closed
}

// Log writes the message to each logger appended at the given level or higher.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Log(level Level, v ...interface{}) {
	if l.Writable() {
		message := l.Formatter.Format(level, v...)
		for _, logger := range l.Loggers.FindByLevel(level) {
			logger.Print(message)
		}

		if l.FatalOn&level > 0 {
			os.Exit(1)
		} else if l.PanicOn&level > 0 {
			panic(message)
		}
	}
}

// Log writes the message to each logger appended at the given level or higher.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Logf(level Level, format string, v ...interface{}) {
	l.Log(level, fmt.Sprintf(format, v...))
}

// Debug prints to each log file at the Debug level.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Debug(v ...interface{}) {
	l.Log(Debug, v...)
}

// Debugf prints to each log file at the Debug level.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.Logf(Debug, format, v...)
}

// Info prints to each log file at the Info level.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Info(v ...interface{}) {
	l.Log(Info, v...)
}

// Infof prints to each log file at the Info level.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Infof(format string, v ...interface{}) {
	l.Logf(Info, format, v...)
}

// Notice prints to each log file at the Notice level.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Notice(v ...interface{}) {
	l.Log(Notice, v...)
}

// Noticef prints to each log file at the Notice level.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Noticef(format string, v ...interface{}) {
	l.Logf(Notice, format, v...)
}

// Warning prints to each log file at the Warning level.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Warning(v ...interface{}) {
	l.Log(Warning, v...)
}

// Warningf prints to each log file at the Warning level.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Warningf(format string, v ...interface{}) {
	l.Logf(Warning, format, v...)
}

// Error prints to each log file at the Error level.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Error(v ...interface{}) {
	l.Log(Error, v...)
}

// Errorf prints to each log file at the Error level.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Logf(Error, format, v...)
}

// Critical prints to each log file at the Critical level.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Critical(v ...interface{}) {
	l.Log(Critical, v...)
}

// Criticalf prints to each log file at the Critical level.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Criticalf(format string, v ...interface{}) {
	l.Logf(Critical, format, v...)
}

// Alert prints to each log file at the Alert level.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Alert(v ...interface{}) {
	l.Log(Alert, v...)
}

// Alertf prints to each log file at the Alert level.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Alertf(format string, v ...interface{}) {
	l.Logf(Alert, format, v...)
}

// Emergency prints to each log file at the Emergency level.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Emergency(v ...interface{}) {
	l.Log(Emergency, v...)
}

// Emergencyf prints to each log file at the Emergency level.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Emergencyf(format string, v ...interface{}) {
	l.Logf(Emergency, format, v...)
}

// open returns a file that logs can be written to.
func (l *Logger) open(name string) *os.File {
	w, err := os.OpenFile(name, FileFlags, FileMode)
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
