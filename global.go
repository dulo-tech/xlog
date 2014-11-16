package xlog

import "io"

// globalInstance stores the global logger.
var globalInstance *Logger

// Close releases any resources held by the global logger. The logger should
// not be used again after calling this method without re-configuring it, as
// this method sets the global instance to nil.
func Close() {
	instance().Close()
	globalInstance = nil
}

// SetName sets the name of the global logger.
func SetName(name string) {
	instance().SetName(name)
}

// Append adds a file to the global logger.
func Append(file string, level Level) {
	instance().Append(file, level)
}

// MultiAppend adds one or more files to the global logger.
func MultiAppend(files []string, level Level) {
	instance().MultiAppend(files, level)
}

// AppendWriter adds a writer to the global logger.
func AppendWriter(writer io.Writer, level Level) {
	instance().AppendWriter(writer, level)
}

// MultiAppendWriter adds one or more io.Writer instances to the global logger.
func MultiAppendWriter(writers []io.Writer, level Level) {
	instance().MultiAppendWriter(writers, level)
}

// Writable returns true when global logging is enabled, and the global logger
// hasn't been closed.
func Writable() bool {
	return instance().Writable()
}

// Log writes the message to each logger appended to the global logger at given level
// or higher.
func Log(level Level, v ...interface{}) {
	instance().Log(level, v...)
}

// Logf writes the message to each logger appended to the global logger given
// level or higher.
func Logf(level Level, format string, v ...interface{}) {
	instance().Logf(level, format, v...)
}

// Debug writes to the global logger at DebugLevel.
// Arguments are handled in the manner of fmt.Print.
func Debug(v ...interface{}) {
	instance().Debug(v...)
}

// Debugf writes to the global logger at DebugLevel.
// Arguments are handled in the manner of fmt.Printf.
func Debugf(format string, v ...interface{}) {
	instance().Debugf(format, v...)
}

// Info writes to the global logger at InfoLevel.
// Arguments are handled in the manner of fmt.Print.
func Info(v ...interface{}) {
	instance().Info(v...)
}

// Infof writes to the global logger at InfoLevel.
// Arguments are handled in the manner of fmt.Printf.
func Infof(format string, v ...interface{}) {
	instance().Infof(format, v...)
}

// Notice writes to the global logger at NoticeLevel.
// Arguments are handled in the manner of fmt.Print.
func Notice(v ...interface{}) {
	instance().Notice(v...)
}

// Noticef writes to the global logger at NoticeLevel.
// Arguments are handled in the manner of fmt.Printf.
func Noticef(format string, v ...interface{}) {
	instance().Noticef(format, v...)
}

// Warning writes to the global logger at WarningLevel.
// Arguments are handled in the manner of fmt.Print.
func Warning(v ...interface{}) {
	instance().Warning(v...)
}

// Warningf writes to the global logger at WarningLevel.
// Arguments are handled in the manner of fmt.Printf.
func Warningf(format string, v ...interface{}) {
	instance().Warningf(format, v...)
}

// Error writes to the global logger at ErrorLevel.
// Arguments are handled in the manner of fmt.Print.
func Error(v ...interface{}) {
	instance().Error(v...)
}

// Errorf writes to the global logger at ErrorLevel.
// Arguments are handled in the manner of fmt.Printf.
func Errorf(format string, v ...interface{}) {
	instance().Errorf(format, v...)
}

// Critical writes to the global logger at CriticalLevel.
// Arguments are handled in the manner of fmt.Print.
func Critical(v ...interface{}) {
	instance().Critical(v...)
}

// Criticalf writes to the global logger at CriticalLevel.
// Arguments are handled in the manner of fmt.Printf.
func Criticalf(format string, v ...interface{}) {
	instance().Criticalf(format, v...)
}

// Alert writes to the global logger at AlertLevel.
// Arguments are handled in the manner of fmt.Print.
func Alert(v ...interface{}) {
	instance().Alert(v...)
}

// Alertf writes to the global logger at AlertLevel.
// Arguments are handled in the manner of fmt.Printf.
func Alertf(format string, v ...interface{}) {
	instance().Alertf(format, v...)
}

// Emergency writes to the global logger at EmergencyLevel.
// Arguments are handled in the manner of fmt.Print.
func Emergency(v ...interface{}) {
	instance().Emergency(v...)
}

// Emergencyf writes to the global logger at EmergencyLevel.
// Arguments are handled in the manner of fmt.Printf.
func Emergencyf(format string, v ...interface{}) {
	instance().Emergencyf(format, v...)
}

// instance calls panic() when the global logger has not been configured.
func instance() *Logger {
	if globalInstance == nil {
		globalInstance = NewLogger("xlog")
	}
	return globalInstance
}
