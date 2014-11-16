package xlog

import (
	"log"
	"os"
	"fmt"
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

// Aliases maps file aliases to real file pointers.
var Aliases = map[string]*os.File{
	"stdout": os.Stdout,
	"stdin": os.Stdin,
	"stderr": os.Stderr,
}

var (
	// DefaultDateFormat is the date format to use when none has been specified.
	DefaultDateFormat string = "2006-01-02 15:04:05.000"
	
	// DefaultMessageFormat is the message format to use when none has been specified.
	DefaultMessageFormat string = "{date|2006-01-02 15:04:05.000} {name}.{level} {message}"
	
	// DefaultAppendFiles stores the names of files appended to the logger by default.
	DefaultAppendFiles []string
	
	// DefaultAppendWriters stores writers that are appended to the logger by default.
	DefaultAppendWriters []io.Writer
	
	// DefaultAppendLevel is the default level used when appending files from
	// DefaultAppendFiles and DefaultAppendWriters.
	DefaultAppendLevel Level = DebugLevel

	// FileFlags defines the file open options.
	FileOpenFlags int = os.O_RDWR|os.O_CREATE|os.O_APPEND

	// FileMode defines the mode files are opened in.
	FileOpenMode os.FileMode = 0666

	// PanicOnFileErrors defines whether the logger should panic when opening a file
	// fails. When set to false, any file open errors are ignored, and the file won't be
	// appended.
	PanicOnFileErrors = true

	// InitialLoggerCapacity defines the initial capacity for each type of logger.
	InitialLoggerCapacity = 4

	// FatalOn represents levels that causes the application to exit.
	FatalOn Level = 0

	// PanicOn represents levels that causes the application to panic.
	PanicOn Level = 0
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

	// pointers contains any files that have been opened for logging.
	pointers []*os.File

	// closed defines whether the logger has been closed.
	closed bool
}

// NewLogger returns a *Logger instance that's been initialized with default values.
func NewLogger(name string) *Logger {
	logger := &Logger{
		Enabled: true,
		Formatter: NewDefaultFormatter(DefaultMessageFormat, name),
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

// NewMultiLogger returns a *Logger instance that's been initialized with one or
// more files at the given level.
func NewMultiLogger(name string, files []string, level Level) *Logger {
	logger := NewLogger(name)
	logger.MultiAppend(files, level);
	return logger;
}

// NewMultiWriterLogger returns a *Logger instance that's been initialized with one or
// more writers at the given level.
func NewMultiWriterLogger(name string, writers []io.Writer, level Level) *Logger {
	logger := NewLogger(name)
	logger.MultiAppendWriters(writers, level);
	return logger;
}

// NewFormattedLogger returns a *Logger instance using the provided formatter.
func NewFormattedLogger(formatter Formatter) *Logger {
	logger := NewLogger("")
	logger.Formatter = formatter
	return logger
}

// SetName sets the name of the logger.
func (l *Logger) SetName(name string) {
	l.Formatter.SetName(name)
}

// Append adds a file that will be written to at the given level or greater.
// The file argument may be either the full path to a system file, or one of the
// aliases "stdout", "stdin", or "stderr".
func (l *Logger) Append(file string, level Level) {
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
func (l *Logger) MultiAppend(files []string, level Level) {
	for _, file := range files {
		l.Append(file, level)
	}
}

// AppendWriter adds a writer that will be written to at the given level or greater.
func (l *Logger) AppendWriter(writer io.Writer, level Level) {
	l.Loggers.Append(newLogger(writer), level)
}

// MultiAppendWriters adds one or more io.Writer instances to the logger.
func (l *Logger) MultiAppendWriters(writers []io.Writer, level Level) {
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
func (l *Logger) Logf(level Level, format string, v ...interface{}) {
	l.Log(level, fmt.Sprintf(format, v...))
}

// Debug writes to the logger at DebugLevel.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Debug(v ...interface{}) {
	l.Log(DebugLevel, v...)
}

// Debugf writes to the logger at DebugLevel.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.Logf(DebugLevel, format, v...)
}

// Info writes to the logger at InfoLevel.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Info(v ...interface{}) {
	l.Log(InfoLevel, v...)
}

// Infof writes to the logger at InfoLevel.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Infof(format string, v ...interface{}) {
	l.Logf(InfoLevel, format, v...)
}

// Notice writes to the logger at NoticeLevel.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Notice(v ...interface{}) {
	l.Log(NoticeLevel, v...)
}

// Noticef writes to the logger at NoticeLevel.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Noticef(format string, v ...interface{}) {
	l.Logf(NoticeLevel, format, v...)
}

// Warning writes to the logger at WarningLevel.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Warning(v ...interface{}) {
	l.Log(WarningLevel, v...)
}

// Warningf writes to the logger at WarningLevel.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Warningf(format string, v ...interface{}) {
	l.Logf(WarningLevel, format, v...)
}

// Error writes to the logger at ErrorLevel.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Error(v ...interface{}) {
	l.Log(ErrorLevel, v...)
}

// Errorf writes to the logger at ErrorLevel.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Logf(ErrorLevel, format, v...)
}

// Critical writes to the logger at CriticalLevel.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Critical(v ...interface{}) {
	l.Log(CriticalLevel, v...)
}

// Criticalf writes to the logger at CriticalLevel.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Criticalf(format string, v ...interface{}) {
	l.Logf(CriticalLevel, format, v...)
}

// Alert writes to the logger at AlertLevel.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Alert(v ...interface{}) {
	l.Log(AlertLevel, v...)
}

// Alertf writes to the logger at AlertLevel.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Alertf(format string, v ...interface{}) {
	l.Logf(AlertLevel, format, v...)
}

// Emergency writes to the logger at EmergencyLevel.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Emergency(v ...interface{}) {
	l.Log(EmergencyLevel, v...)
}

// Emergencyf writes to the logger at EmergencyLevel.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Emergencyf(format string, v ...interface{}) {
	l.Logf(EmergencyLevel, format, v...)
}

// open returns a file that logs can be written to.
func (l *Logger) open(name string) *os.File {
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

