package xlog

import "io"

var (
	// globalInstance stores the global logger.
	globalInstance *DefaultLogger

	// globalAppended stores whether files have been appended to the global logger.
	globalAppended bool = false

	// globalLoggers stores the loggers created by the GetLogger() function.
	globalLoggers map[string]*DefaultLogger
)

// Instance returns the global logger.
func Instance() *DefaultLogger {
	if globalInstance == nil {
		globalInstance = New("xlog")
		globalInstance.Append("stdout", DebugLevel)
		globalAppended = false
	}

	return globalInstance
}

// GetLogger returns the *DefaultLogger with the given name. The logger will be
// created if it's not already been created. Only a single *DefaultLogger instance
// is created for a name.
func GetLogger(name string) *DefaultLogger {
	if globalLoggers == nil {
		globalLoggers = make(map[string]*DefaultLogger)
	}
	if _, ok := globalLoggers[name]; !ok {
		globalLoggers[name] = New(name)
	}

	return globalLoggers[name]
}

// Close releases any resources held by the global logger. The logger should
// not be used again after calling this method without re-configuring it, as
// this method sets the global instance to nil.
func Close() {
	Instance().Close()
	globalInstance = nil
}

// SetName sets the name of the global logger.
func SetName(name string) {
	Instance().Name = name
}

// SetFormatter sets the formatter used by the global logger.
func SetFormatter(formatter Formatter) {
	Instance().Formatter = formatter
}

// SetContainer sets the logger container used by the global logger.
func SetContainer(lc Container) {
	Instance().Container = lc
}

// Enabled returns whether the global logger is enabled.
func Enabled() bool {
	return Instance().Enabled
}

// SetEnabled sets whether the global logger is enabled.
func SetEnabled(enabled bool) {
	Instance().Enabled = enabled
}

// Append adds a file to the global logger.
func Append(file string, level Level) {
	if !globalAppended {
		Instance().ClearAppended()
		globalAppended = true
	}
	Instance().Append(file, level)
}

// MultiAppend adds one or more files to the global logger.
func MultiAppend(files []string, level Level) {
	if !globalAppended {
		Instance().ClearAppended()
		globalAppended = true
	}
	Instance().MultiAppend(files, level)
}

// AppendWriter adds a writer to the global logger.
func AppendWriter(writer io.Writer, level Level) {
	if !globalAppended {
		Instance().ClearAppended()
		globalAppended = true
	}
	Instance().AppendWriter(writer, level)
}

// MultiAppendWriters adds one or more io.Writer instances to the global logger.
func MultiAppendWriters(writers []io.Writer, level Level) {
	if !globalAppended {
		Instance().ClearAppended()
		globalAppended = true
	}
	Instance().MultiAppendWriters(writers, level)
}

// Writable returns true when global logging is enabled, and the global logger
// hasn't been closed.
func Writable() bool {
	return Instance().Writable()
}

// Log writes the message to each logger appended to the global logger at given level
// or higher.
func Log(level Level, v ...interface{}) {
	Instance().Log(level, v...)
}

// Logf writes the message to each logger appended to the global logger given
// level or higher.
func Logf(level Level, format string, v ...interface{}) {
	Instance().Logf(level, format, v...)
}

// Debug writes to the global logger at DebugLevel.
// Arguments are handled in the manner of fmt.Print.
func Debug(v ...interface{}) {
	Instance().Debug(v...)
}

// Debugf writes to the global logger at DebugLevel.
// Arguments are handled in the manner of fmt.Printf.
func Debugf(format string, v ...interface{}) {
	Instance().Debugf(format, v...)
}

// Info writes to the global logger at InfoLevel.
// Arguments are handled in the manner of fmt.Print.
func Info(v ...interface{}) {
	Instance().Info(v...)
}

// Infof writes to the global logger at InfoLevel.
// Arguments are handled in the manner of fmt.Printf.
func Infof(format string, v ...interface{}) {
	Instance().Infof(format, v...)
}

// Notice writes to the global logger at NoticeLevel.
// Arguments are handled in the manner of fmt.Print.
func Notice(v ...interface{}) {
	Instance().Notice(v...)
}

// Noticef writes to the global logger at NoticeLevel.
// Arguments are handled in the manner of fmt.Printf.
func Noticef(format string, v ...interface{}) {
	Instance().Noticef(format, v...)
}

// Warning writes to the global logger at WarningLevel.
// Arguments are handled in the manner of fmt.Print.
func Warning(v ...interface{}) {
	Instance().Warning(v...)
}

// Warningf writes to the global logger at WarningLevel.
// Arguments are handled in the manner of fmt.Printf.
func Warningf(format string, v ...interface{}) {
	Instance().Warningf(format, v...)
}

// Error writes to the global logger at ErrorLevel.
// Arguments are handled in the manner of fmt.Print.
func Error(v ...interface{}) {
	Instance().Error(v...)
}

// Errorf writes to the global logger at ErrorLevel.
// Arguments are handled in the manner of fmt.Printf.
func Errorf(format string, v ...interface{}) {
	Instance().Errorf(format, v...)
}

// Critical writes to the global logger at CriticalLevel.
// Arguments are handled in the manner of fmt.Print.
func Critical(v ...interface{}) {
	Instance().Critical(v...)
}

// Criticalf writes to the global logger at CriticalLevel.
// Arguments are handled in the manner of fmt.Printf.
func Criticalf(format string, v ...interface{}) {
	Instance().Criticalf(format, v...)
}

// Alert writes to the global logger at AlertLevel.
// Arguments are handled in the manner of fmt.Print.
func Alert(v ...interface{}) {
	Instance().Alert(v...)
}

// Alertf writes to the global logger at AlertLevel.
// Arguments are handled in the manner of fmt.Printf.
func Alertf(format string, v ...interface{}) {
	Instance().Alertf(format, v...)
}

// Emergency writes to the global logger at EmergencyLevel.
// Arguments are handled in the manner of fmt.Print.
func Emergency(v ...interface{}) {
	Instance().Emergency(v...)
}

// Emergencyf writes to the global logger at EmergencyLevel.
// Arguments are handled in the manner of fmt.Printf.
func Emergencyf(format string, v ...interface{}) {
	Instance().Emergencyf(format, v...)
}
