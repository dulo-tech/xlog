package xlog

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Formatter is an interface that provides methods that format log messages.
type Formatter interface {
	SetFormat(format string)
	PlaceholderFunc(key string, f func(string) string)
	Format(name string, level Level, v ...interface{}) string
}

// DefaultFormatter is the default implementation of the Formatter interface.
type DefaultFormatter struct {
	messageFormat string
	dateFormat    string
	funcs         map[string]func(string) string
}

// NewDefaultFormatter creates and returns a new DefaultFormatter instance.
func NewDefaultFormatter(messageFormat, dateFormat string) *DefaultFormatter {
	messageFormat, dateFormat = SanitizeForDate(messageFormat, dateFormat)
	placeholders := make(map[string]func(string) string)
	return &DefaultFormatter{messageFormat, dateFormat, placeholders}
}

// SetFormat changes the set message format.
func (f *DefaultFormatter) SetFormat(format string) {
	f.messageFormat, f.dateFormat = SanitizeForDate(format, f.dateFormat)
}

// PlaceholderFunc adds a callback function which provides a replacement for key in a string format.
func (f *DefaultFormatter) PlaceholderFunc(key string, fn func(string) string) {
	f.funcs[key] = fn
}

// Format formats a log message for the given level.
func (f *DefaultFormatter) Format(name string, level Level, v ...interface{}) string {
	placeholders := map[string]string{
		"{date}":    (time.Now()).Format(f.dateFormat),
		"{level}":   Levels[level],
		"{message}": fmt.Sprint(v...),
		"{name}":    name,
	}

	formatted := f.messageFormat
	for placeholder, value := range placeholders {
		formatted = strings.Replace(formatted, placeholder, value, -1)
	}
	for key, fn := range f.funcs {
		formatted = strings.Replace(formatted, "{"+key+"}", fn(key), -1)
	}

	return formatted
}

// SanitizeForDate replaces date placeholders containing a date format with
// a plain {date} placeholder. The altered message format is returned, along
// with the found date format.
func SanitizeForDate(messageFormat, dateFormat string) (string, string) {
	regex := regexp.MustCompile(`{date\|([^}]+)}`)
	captured := regex.FindStringSubmatch(messageFormat)
	if len(captured) == 2 {
		dateFormat = captured[1]
		messageFormat = regex.ReplaceAllString(messageFormat, "{date}")
	}

	return messageFormat, dateFormat
}
