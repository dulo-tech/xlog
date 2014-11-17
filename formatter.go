package xlog

import (
	"fmt"
	"strings"
	"time"
	"regexp"
)

// Formatter is an interface that provides methods that format log messages.
type Formatter interface {
	Format(name string, level Level, v ...interface{}) string
}

// DefaultFormatter is the default implementation of the Formatter interface.
type DefaultFormatter struct {
	messageFormat string
	dateFormat    string
}

// NewDefaultFormatter creates and returns a new DefaultFormatter instance.
func NewDefaultFormatter(messageFormat string) *DefaultFormatter {
	messageFormat, dateFormat := SanitizeForDate(messageFormat)
	return &DefaultFormatter{messageFormat, dateFormat}
}

// SetMessageFormat changes the set message format.
func (f *DefaultFormatter) SetMessageFormat(messageFormat string) {
	f.messageFormat, f.dateFormat = SanitizeForDate(messageFormat)
}

// Format formats a log message for the given level.
func (f *DefaultFormatter) Format(name string, level Level, v ...interface{}) string {
	placeholders := map[string]string{
		"{date}": (time.Now()).Format(f.dateFormat),
		"{level}": Levels[level],
		"{message}": fmt.Sprint(v...),
		"{name}": name,
	}

	formatted := f.messageFormat
	for placeholder, value := range placeholders {
		formatted = strings.Replace(formatted, placeholder, value, -1)
	}

	return formatted
}

// SanitizeForDate replaces date placeholders containing a date format with
// a plain {date} placeholder. The altered message format is returned, along
// with the found date format.
func SanitizeForDate(messageFormat string) (string, string) {
	regex := regexp.MustCompile(`{date\|([^}]+)}`)
	dateFormat := DefaultDateFormat
	captured := regex.FindStringSubmatch(messageFormat)
	if len(captured) == 2 {
		dateFormat = captured[1]
		messageFormat = regex.ReplaceAllString(messageFormat, "{date}")
	}

	return messageFormat, dateFormat
}

