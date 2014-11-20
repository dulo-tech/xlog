package xlog

// LoggerWriter wraps a logger in an io.Writer instance.
type LoggerWriter struct {
	// logger is the wrapped logger.
	logger Loggable
	
	// level is the level being written to.
	level Level
}

// NewLoggerWriter returns a new *LoggerWriter instance.
func NewLoggerWriter(logger Loggable, level Level) *LoggerWriter {
	return &LoggerWriter{logger, level}
}

// Write implements io.Writer.Write.
func (w *LoggerWriter) Write(p []byte) (int, error) {
	w.logger.Log(w.level, string(p))
	return len(p), nil
}
