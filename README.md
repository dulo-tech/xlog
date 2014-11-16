xLog
====
A simple logger for the Go language.

Documentation is available from the [GoDoc website](http://godoc.org/github.com/dulo-tech/xlog).

* [Installation](#installation)
* [Examples](#examples)
* [Global Configuration](#global-configuration)
* [Global Logger](#global-logger)
* [Custom Formatters](#custom-formatters)
* [Custom Level Behavior](#custom-level-behavior)
* [Loggable Interface](#loggable-interface)
* [License](#license)


#### Installation
Use the go get command to fetch the package.  
`go get github.com/dulo-tech/xlog`

Then use the import statement in your source code to use the package.
```go
import "github.com/dulo-tech/xlog"
```


#### Examples
```go
package main

import (
    "os"
    "github.com/dulo-tech/xlog"
)

func main() {
    // Start by creating a new logger. Each logger needs a name. The name is
    // usually reflective of the sub-system that will be logging messages. The
    // name will appear in the logged messages so each message can be linked to
    // a specific area of your code.
    logger := xlog.NewLogger("testing")
    
    // Add files where the logs will be written. You can use a full system path
    // to a file, or the aliases "stdout", "stderr", and "stdin". Each file you
    // add is assigned a logging level. The file is written to when a message is
    // logged to that level or a greater level. The levels are:
    //
    // xlog.DebugLevel
    // xlog.InfoLevel
    // xlog.NoticeLevel
    // xlog.WarningLevel
    // xlog.ErrorLevel
    // xlog.CriticalLevel
    // xlog.AlertLevel
    // xlog.EmergencyLevel
    //
    // When you append a file using the xlog.DebugLevel level, the file will be
    // written to for every level, because every level is greater than xdt.Debug.
    logger.Append("stdout", xlog.DebugLevel)
    
    // For each log level there is a xdt.Loggable method used to write messages
    // to that level. For example xdt.Loggable.Debug(), xdt.Loggable.Alert(), etc.
    // There is a formatting method for each level as well. For example
    // xdt.Loggable.Debugf(), xdt.Loggable.Alertf(), etc.
    
    // Outputs: 2014-11-15 09:40:28.693 testing.DEBUG Test debug message.
    logger.Debug("Test debug message.")
    
    // Outputs: 2014-11-15 09:40:28.693 testing.DEBUG Test debug message.
    logger.Debugf("Test %s message.", "debug")
    
    // Outputs: 2014-11-15 09:40:28.701 testing.WARNING Test warning message.
    logger.Warning("Test warning message.")
    
    // Outputs: 2014-11-15 09:40:28.723 testing.INFO Test info message.
    logger.Infof("Test %s message.", "info")
    
    // Any log message with the xlog.WarningLevel level and above will be logged to
    // stdout.
    logger = xlog.NewLogger("testing")
    logger.Append("stdout", xlog.WarningLevel)
    
    // This doesn't output anything because the xlog.DebugLevel level is lower than
    // the xlog.WarningLevel level.
    logger.Debug("Test debug message.")
    
    // Outputs: 2014-11-15 09:40:28.701 testing.WARNING Test warning message.
    logger.Warning("Test warning message.")
    
    // This doesn't output anything because the Info level is lower than the
    // Warning level.
    logger.Infof("Test %s message.", "info")
    
    // Logging can be disabled, and when disabled the calls to the logging
    // methods will simply be ignored.
    
    // Outputs: 2014-11-15 09:40:28.701 testing.NOTICE Test notice message.
    logger.Notice("Test notice message.")
    
    // Now disable logging.
    logger.Enabled = false
    
    // This doesn't output anything because logging is now disabled, but you
    // can still call the methods without any errors.
    logger.Notice("Test notice message.")
    
    // You can append as many files as you want. This logs messages to the
    // xlog.DebugLevel and above to stdout, and messages xlog.ErrorLevel level
    // and above to a file.
    // Note, when you have the logger open files, the files need to be closed by
    // the logger by deferring the defer logger.Close() method. The logger cannot
    // be used once it's been closed.
    logger = xlog.NewLogger("testing")
    logger.Append("stdout", xlog.DebugLevel)
    logger.Append("/var/logs/main-error.log", xlog.ErrorLevel)
    defer logger.Close()
    
    // You can manage the files yourself by using the logger.AppendWriter()
    // method.
    fp, err := os.OpenFile(
        "/var/logs/main-error.log",
        os.O_RDWR|os.O_CREATE | os.O_APPEND,
        0666,
    )
    if err != nil {
        panic(err)
    }
    defer fp.Close()
    logger.AppendWriter(fp, xlog.DebugLevel)
    logger.AppendWriter(os.Stderr, xlog.ErrorLevel)
    
    // You can append multiple files using xdt.Logger.MultiAppend() and
    // xdt.Logger.MultiAppendWriter().
    logger = xlog.NewLogger("testing")
	logger.MultiAppend(
		[]string{
			"stderr",
			"/var/logs/main-error.log",
		},
		xlog.DebugLevel,
	)
    defer logger.Close()
    
    // You can also create the logger with an array of file names or writers.
	files := []string{
        "stderr",
        "/var/logs/main-error.log",
	}
	logger = xlog.NewMultiLogger("testing", files, xdt.Debug)
    
    // Change the way the log messages are formatted. The xlog.Formatter interface
    // requires a format string and a name. The format string defines how the
    // log messages are formatted. The name is the same value you would pass to
    // a new xdt.Logger instance. Several placeholders can be used in the format
    // string, which will be replaced by actual values:
    //
    // {date} The datetime when the message was logged.
    // {name} The name of the logger/formatter.
    // {level} A string representation of the log level.
    // {message} The message that was logged.
    logger.Formatter = xlog.NewDefaultFormatter(
        "{date} {name} - {level} - {message}",
        "testing",
    )
    
    // Outputs: 2014-11-15 09:54:16.278 testing - DEBUG - This is a debug test.
    logger.Debug("Test debug message.")
    
    // Change the way dates are printed using the Go time syntax inside the
    // {date} placeholder.
    // See: http://golang.org/pkg/time/#Time.Format
    logger.Formatter = xlog.NewDefaultFormatter(
        "{date|Jan _2 15:04:05} {level} {message}",
        "testing",
    )
    
    // Outputs: Nov 15 09:56:56 DEBUG Test debug message.
    logger.Debug("Test debug message.")
    
    // Creating a logger with a pre-configured formatter.
    formatter := xlog.NewDefaultFormatter(
        "{date} {name} [{level}] {message}",
        "testing",
    )
    logger = xlog.NewFormattedLogger(formatter)
    logger.Append("stdout", xlog.DebugLevel)
    
    // Outputs: 2014-11-15 09:59:32.427 testing [DEBUG] Test debug message.
    logger.Debug("Test debug message.")
    
    // The message format can be changed without setting a new Formatter. The
    // name can be changed as well.
    logger.Formatter.SetMessageFormat("{date} {message}")
    logger.Formatter.SetName("debug-testing")
}
```


#### Global Configuration
The follow examples demonstrate the use of the xlog global configuration values.
Changing these values effects every logger.

```go
package main

import (
    "io"
    "os"
    "time"
    "github.com/dulo-tech/xlog"
)

func main() {
    // The first thing you need to do after creating a new logger is append
    // one or more files to it. That must be done each and every time you
    // create a logger. The alternative is to set the names of the files
    // globally. Once done the files will be automatically appended to each
    // new logger. You can also set global writers, which will also be
    // automatically appended.
    xlog.DefaultAppendFiles = []string{"stdout", "/var/log/messages.log"}
    xlog.DefaultAppendWriters = []io.Writer{os.Stdout, os.Stderr}
    
    // The files and writers that you have automatically appended will by
    // default be appended at the xlog.DebugLevel, but that can also be changed.
    xlog.DefaultAppendLevel = xlog.WarningLevel

    // Change the message format for each new logger.
    xlog.DefaultMessageFormat = "{date} {message}"
    xlog.DefaultMessageFormat = "{date|2006-01-02 15:04:05.000} {level} {message}"
    
    // Change the date format for each new logger. Either write it yourself,
    // or use one of the defaults from the time package. Note that changing
    // this global value has no effect when the message format string already
    // specifies a date format, eg "{date|2006-01-02}". The default date
    // format only applies to the "{date}" placeholder.
    xlog.DefaultDateFormat = "2006-01-02 15:04:05.000"
    xlog.DefaultDateFormat = time.UnixDate
    xlog.DefaultDateFormat = time.StampMicro

    // You can replicate the functionality of Go's system logger log.Fatal()
    // and log.Panic() using logger.FatalOn and logger.PanicOn.
    
    // Logging a message to xlog.CriticalLevel will cause a fatal shut down
    // using os.Exit(1).
    xlog.FatalOn = xlog.CriticalLevel
    
    // Logging a message to either xlog.AlertLevel or xlog.EmergencyLevel
     // causes a panic using panic().
    xlog.PanicOn = xlog.AlertLevel | xlog.EmergencyLevel
    
    // Change the mode and permissions used when the logger opens a file.
    xlog.FileOpenFlags = os.O_RDWR|os.O_CREATE|os.O_APPEND
    xlog.FileOpenMode = 0666
    
    // The logger will panic by calling panic() when it fails to open a file.
    // The panic can be globally suppressed. When the panics are suppressed,
    // and the logger fails to open a file, the file will simply be ignored.
    // No logs will be written to it.
    xlog.PanicOnFileErrors = false
    
    // You can increase the loggers initial capacity for appended files, which
    // may help with performance when you know the loggers being created will
    // have more than 4 (the default) files appended. The logger uses this value
    // with the make() function when allocating internal maps.
    xlog.InitialLoggerCapacity = 10
    
    // Each log level has a corresponding string representation which is used
    // in the log messages. Those can be changed. Here we change the string
    // representations of xlog.DebugLevel and xlog.InfoLevel from their default
    // values ("DEBUG", "INFO") to "Debug" and "Info".
    xlog.Levels[xlog.DebugLevel] = "Debug"
    xlog.Levels[xlog.InfoLevel] = "Info"
    
    // The strings "stdout", "stderr", and "stdin" may be passed as a file name
    // to the logger append methods, which is useful when the files to be written
    // to are saved as strings in a configuration file, or passed as strings at
    // the command line. Be default the aliases map to os.Stdout, os.Stderr, and
    // os.Stdin, but those can be changed to any writer.
    fp, err := os.OpenFile(
        "/var/logs/output.log",
        os.O_RDWR|os.O_CREATE | os.O_APPEND,
        0666,
    )
    if err != nil {
        panic(err)
    }
    defer fp.Close()
    
    xlog.Aliases["stdout"] = fp
    xlog.Aliases["stderr"] = fp
    
    // You can even create your own aliases through the xlog.Aliases variable,
    // and then append the file using the alias.
    xlog.Aliases["output"] = fp
    logger := NewLogger()
    logger.Append("output", xlog.DebugLevel)
}
```


#### Global Logger
You can use the global logger when your application is small, and does not need
to log messages from several sub-systems using a different names. The global
logger is easy to use, but it's not quite as flexible as using individual logger
instances.

```go
package main

import (
    "os"
    "fmt"
    "github.com/dulo-tech/xlog"
)

func main() {
    // The global logger has the same methods as a logger
    // instance (Debug(), Warningf(), Info(), etc), but as exposed functions
    // in the package scope.
    
    // Outputs: 2014-11-15 09:40:28.693 xlog.DEBUG Test debug message.
    xlog.Debug("Test debug message.")
    
    // Outputs: 2014-11-15 09:40:28.693 xlog.DEBUG Test debug message.
    xlog.Debugf("Test %s message.", "debug")
    
    // Outputs: 2014-11-15 09:40:28.701 xlog.WARNING Test warning message.
    xlog.Warning("Test warning message.")

    // By default the global logger has the name "xlog", but that can be
    // changed.
    xlog.SetName("testing")
    
    // Be default the global logger writes messages at xlog.DebugLevel and
    // above to stdout. That can be changed by appending the files you want
    // at the level you want. Note that calling any of the global Append
    // functions removes the default stdout. You'll have to add it back
    // if you still want to log to stdout.
    xlog.Append("stderr", xlog.WarningLevel)
    xlog.Append("stdout", xlog.DebugLevel)
    xlog.AppendWriter(os.Stdout, xlog.InfoLevel)
    
    // Outputs: 2014-11-15 09:40:28.693 testing.DEBUG Test debug message.
    xlog.Debug("Test debug message.")
    
    // Outputs: 2014-11-15 09:40:28.701 testing.WARNING Test warning message.
    xlog.Warning("Test warning message.")
    
    // Just like the logger instances, you need to ensure the global logger
    // closes any files it opened. You can continue to use the logger after
    // closing it, but closing the logger resets it's configuration, which
    // defaults the global logger to the default name and stdout.
    xlog.Append("/var/log/messages.log", xlog.DebugLevel)
    defer xlog.Close()
    
    // You can disable/enable global logging and test to see if it's enabled.
    xlog.SetEnabled(false)
    if !xlog.Enabled() {
        fmt.Println("Logging is disabled.")
    }
}
```


#### Custom Formatters
You can create your own message formatter by creating a struct that implements
the `xlog.Formatter` interface, which has the following signature:

```go
type Formatter interface {
    // Name returns the name of the formatter.
    Name() string
    
    // SetName sets the name of the formatter.
	SetName(name string)
	
	// Format formats a log message for the given level.
	Format(level Level, v ...interface{}) string
}
```

This example creates a formatter than always formats messages into an empty
string. The logger discards empty messages, which means this formatter causes
all messages to be discarded.


```go
package main

import "github.com/dulo-tech/xlog"

// NullFormatter implements the xlog.Formatter interface where all
// log messages are discarded.
type NullFormatter struct {
	name string
}

// Name returns the name of the formatter.
func (f *NullFormatter) Name() string {
    return f.name
}

// SetName sets the name of the formatter.
func (f *NullFormatter) SetName(name string) {
	f.name = name
}

// Format formats a log message for the given level.
func (f *NullFormatter) Format(level Level, v ...interface{}) string {
    return ""
}

func main() {
    // Creating a logger which discards all messages.
    formatter := &NullFormatter{""}
    logger = xlog.NewFormattedLogger(formatter)
    logger.Append("stdout", xlog.DebugLevel)
}
```

Internally the `xlog` package uses the standard Go logger, `log.Logger`. An instance
of `log.Logger` is created for each file you append to the logger. The `log.Logger`
instances are managed by the `xlog.LoggerMap` interface, which stores the loggers
and makes them retrievable by level. The `xlog.LoggerMap` interface has the
following signature:

```go
type LoggerMap interface {
    // Append adds a logger to the map for the given level.
	Append(logger *log.Logger, level Level)
	
	// FindByLevel returns all the loggers added to the map at the given level
	// or greater.
	FindByLevel(level Level) []*log.Logger
}
```

#### Custom Level Behavior
By default when you log a message to `xlog.DebugLevel`, the message is written
to all files added at the `xlog.DebugLevel` level *and greater*. The
`xlog.LoggerMap.FindByLevel()` method is responsible for returning loggers registered
at a given level and all those registered at greater levels.

If you wanted logs written at a given level to only be written at that level, and
not levels greater than it, you can implement your own logger map which changes
the default behavior.


```go
package main

import (
    "log"
    "github.com/dulo-tech/xlog"
)

// CustomLoggerMap maps loggers to levels.
type CustomLoggerMap struct {
	loggers map[xlog.Level][]*log.Logger
}

// NewCustomLoggerMap creates and returns a *CustomLoggerMap instance.
func NewCustomLoggerMap() *CustomLoggerMap {
    // Make the internal map the size of xlog.Levels, and initialize
    // each slice to the default initial capacity.
	loggers := make(map[xlog.Level][]*log.Logger, len(xlog.Levels))
	for level, _ := range xlog.Levels {
		loggers[level] = make([]*log.Logger, 0, xlog.InitialLoggerCapacity)
	}

	return &CustomLoggerMap{loggers}
}

// Append adds a logger to the map at the given level.
func (m *CustomLoggerMap) Append(logger *log.Logger, level xlog.Level) {
    m.loggers[level] = append(m.loggers[level], logger)
}

// FindByLevel returns the loggers at the given level, and only the given level.
func (m *CustomLoggerMap) FindByLevel(level xlog.Level) []*log.Logger {
	return m.loggers[level]
}

func main() {
    // Creating a logger that uses the custom logger map.
    logger = xlog.NewLogger("testing")
    logger.Loggers = NewCustomLoggerMap()
    logger.Append("stdout", xlog.DebugLevel)
}
```

#### Loggable Interface
The `xlog.NewLogger()` method and other New methods return an instance of
the struct `xlog.Logger`, which implements the `xlog.Loggable` interface.

```go
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
```

Note the absence of methods like `Append()`, `MultiAppend()`, `AppendWriter()`,
and `Close()`, which are members of the `xlog.Logger` struct. The `xlog.Loggable`
interface does not concern itself for how a logger is configured or it's
lifecycle. The `xlog.Loggable` interface only exposes methods for logging
messages at various log levels.

You are encouraged to reference `xlog.Loggable` instead of `xlog.Logger` in your
code to keep it flexible to future API changes. For example instead of creating
a struct using `xlog.Logger` like this:

```go
type WebScraper struct {
    logger xlog.Logger
}

func NewWebScraper(logger *xlog.Logger) {
    return &WebScraper{logger}
}
```

You should instead use `xlog.Loggable` like this:

```go
type WebScraper struct {
    logger xlog.Loggable
}

func NewWebScraper(logger *xlog.Loggable) {
    return &WebScraper{logger}
}
```

You should configure your loggers in the main package, and let the rest of your
source code deal with the `xlog.Loggable` interface exclusively.


#### License
xLog is has been released under the MIT license, a copy of which is included in
the LICENSE file, which you can find in the source code. You're encouraged to
use the xLog source code in any way you want, for whatever reason you want.
