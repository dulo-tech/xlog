package xlog

import (
	"os"
	"fmt"
	"io"
)

const (
	// DefaultDateFormat is the date format to use when none has been specified.
	DefaultDateFormat string = "2006-01-02 15:04:05.000"

	// DefaultMessageFormat is the message format to use when none has been specified.
	DefaultMessageFormat string = "{date|2006-01-02 15:04:05.000} {name}.{level} {message}"

	// DefaultFileOpenFlags defines the file open options.
	DefaultFileOpenFlags int = os.O_RDWR|os.O_CREATE | os.O_APPEND

	// DefaultFileOpenMode defines the mode files are opened in.
	DefaultFileOpenMode os.FileMode = 0666

	// DefaultPanicOnFileErrors defines whether the logger should panic when opening a file
	// fails. When set to false, any file open errors are ignored, and the file won't be
	// appended.
	DefaultPanicOnFileErrors = true

	// DefaultInitialCapacity defines the initial capacity for each type of logger.
	DefaultInitialCapacity = 4
)

// Settings represents a group of logger settings.
type Settings struct {
	// Enabled defines whether logging is enabled.
	Enabled  bool

	// Formatter is used to format the log messages.
	Formatter

	// Container holds the appended file loggers.
	Container

	// FatalOn represents levels that causes the application to exit.
	FatalOn Level

	// PanicOn represents levels that causes the application to panic.
	PanicOn Level

	// FileFlags defines the file open options.
	FileOpenFlags int

	// FileMode defines the mode files are opened in.
	FileOpenMode os.FileMode
	
	// PanicOnFileErrors defines whether the logger should panic when opening a file
	// fails. When set to false, any file open errors are ignored, and the file won't be
	// appended.
	PanicOnFileErrors bool
}

// NewDefaultSettings returns a new *Settings instance.
func NewDefaultSettings(enabled bool) *Settings {
	return &Settings{
		Enabled: enabled,
		Formatter: NewDefaultFormatter(DefaultMessageFormat, DefaultDateFormat),
		Container: NewDefaultContainer(DefaultInitialCapacity),
		FileOpenFlags: DefaultFileOpenFlags,
		FileOpenMode: DefaultFileOpenMode,
		PanicOnFileErrors: DefaultPanicOnFileErrors,
	}
}

// DefaultLogger is the default implementation of the Loggable interface.
type DefaultLogger struct {
	// Name of the logger.
	Name string
	
	// Settings for the logger.
	*Settings
}

// NewFromSettings returns a *DefaultLogger instance which uses the provided settings.
func NewFromSettings(name string, settings *Settings) *DefaultLogger {
	return &DefaultLogger{
		Name: name,
		Settings: settings,
	}
}

// New returns a *DefaultLogger instance that's been initialized with default values.
func New(name string) *DefaultLogger {
	return NewFromSettings(name, NewDefaultSettings(true))
}

// NewFiles returns a *DefaultLogger instance that's been initialized with one or
// more files at the given level.
func NewFiles(name string, files []string, level Level) *DefaultLogger {
	logger := New(name)
	logger.MultiAppend(files, level);
	return logger;
}

// NewWriters returns a *DefaultLogger instance that's been initialized with one or
// more writers at the given level.
func NewWriters(name string, writers []io.Writer, level Level) *DefaultLogger {
	logger := New(name)
	logger.MultiAppendWriters(writers, level);
	return logger;
}

// Writable returns true when logging is enabled, and the logger hasn't been closed.
func (l *DefaultLogger) Writable() bool {
	return l.Enabled && !l.Settings.Container.Closed()
}

// Closed returns whether the logger has been closed.
func (l *DefaultLogger) Closed() bool {
	return l.Settings.Container.Closed()
}

// Close disables logging and frees up resources used by the logger.
// Note this method only closes files opened by the logger. It's the user's
// responsibility to close files that were passed to the logger via the
// AppendWriter method.
func (l *DefaultLogger) Close() {
	l.Settings.Container.Close()
	l.Settings.Enabled = false
}

// Append adds a file that will be written to at the given level or greater.
// The file argument may be either the full path to a system file, or one of the
// aliases "stdout", "stdin", or "stderr".
func (l *DefaultLogger) Append(file string, level Level) {
	if w, ok := Aliases[file]; ok {
		l.Container.Append(w, level)
	} else {
		w := l.open(file)
		if w != nil {
			l.Container.Append(w, level)
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
	l.Container.Append(writer, level)
}

// MultiAppendWriters adds one or more io.Writer instances to the logger.
func (l *DefaultLogger) MultiAppendWriters(writers []io.Writer, level Level) {
	for _, writer := range writers {
		l.AppendWriter(writer, level)
	}
}

// ClearAppended removes all the files that have been appended to the logger.
func (l *DefaultLogger) ClearAppended() {
	l.Container.Clear()
}

// Log writes the message to each logger appended at the given level or higher.
// Arguments are handled in the manner of fmt.Print.
func (l *DefaultLogger) Log(level Level, v ...interface{}) {
	if l.Writable() {
		message := l.Formatter.Format(l.Name, level, v...)
		if message != "" {
			for _, logger := range l.Container.Get(level) {
				logger.Print(message)
			}
		}

		if l.Settings.FatalOn&level > 0 {
			os.Exit(1)
		} else if l.Settings.PanicOn&level > 0 {
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

// Writer returns a *LoggerWriter instance which wraps this logger.
func (l *DefaultLogger) Writer(level Level) *LoggerWriter {
	return NewLoggerWriter(l, level)
}

// open returns a file that logs can be written to.
func (l *DefaultLogger) open(name string) *os.File {
	w, err := os.OpenFile(name, l.Settings.FileOpenFlags, l.Settings.FileOpenMode)
	if err != nil {
		if l.Settings.PanicOnFileErrors {
			panic(err)
		} else {
			w = nil
		}
	}

	return w
}

