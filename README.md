xLog
====
A simple logger for the Go language.

Documentation is available from the [GoDoc website](http://godoc.org/github.com/dulo-tech/xlog).

* [Installation](#installation)
* [Examples](#examples)
* [Global Logger](#global-logger)
* [Custom Formatters](#custom-formatters)
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
    logger := xlog.New("testing")
    
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
    logger = xlog.New("testing")
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
    logger.Settings.Enabled = false
    
    // This doesn't output anything because logging is now disabled, but you
    // can still call the methods without any errors.
    logger.Notice("Test notice message.")
    
    // You can append as many files as you want. This logs messages to the
    // xlog.DebugLevel and above to stdout, and messages xlog.ErrorLevel level
    // and above to a file.
    // Note, when you have the logger open files, the files need to be closed by
    // the logger by deferring the defer logger.Close() method. The logger cannot
    // be used once it's been closed.
    logger = xlog.New("testing")
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
    logger = xlog.New("testing")
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
	logger = xlog.NewFiles("testing", files, xdt.Debug)
    
    // Change the way the log messages are formatted. The xlog.Formatter interface
    // requires a format string. The format string defines how the
    // log messages are formatted. Several placeholders can be used in the format
    // string, which will be replaced by actual values:
    //
    // {date} The datetime when the message was logged.
    // {name} The name of the logger.
    // {level} A string representation of the log level.
    // {message} The message that was logged.
    logger.Settings.Formatter = xlog.NewDefaultFormatter(
        "{date} {name} - {level} - {message}",
        DefaultDateFormat,
    )
    
    // Outputs: 2014-11-15 09:54:16.278 testing - DEBUG - This is a debug test.
    logger.Debug("Test debug message.")
    
    // Change the way dates are printed using the Go time syntax inside the
    // {date} placeholder.
    // See: http://golang.org/pkg/time/#Time.Format
    logger.Settings.Formatter = xlog.NewDefaultFormatter(
        "{date|Jan _2 15:04:05} {level} {message}",
        DefaultDateFormat
    )
    
    // Outputs: Nov 15 09:56:56 DEBUG Test debug message.
    logger.Debug("Test debug message.")
    
    // Creating a logger with a pre-configured formatter.
    logger = xlog.New("testing")
    logger.Append("stdout", xlog.DebugLevel)
    logger.Settings.formatter = xlog.NewDefaultFormatter(
        "{date} {name} [{level}] {message}",
        DefaultDateFormat
    )
    
    // Outputs: 2014-11-15 09:59:32.427 testing [DEBUG] Test debug message.
    logger.Debug("Test debug message.")
    
    // The message format can be changed without setting a new Formatter.
    logger.Settings.Formatter.SetMessageFormat("{date} {message}")
    
    // In addition to the default placeholders like {date} and {message}, you
    // can also define your own. The Formatter.PlaceholderFunc() takes the value
    // for the placeholder, and a function which returns the value. Below we
    // define the placeholder {hostname} which will be replaced by the OS hostname.
    logger.Settings.Formatter.PlaceholderFunc("hostname", func(key string) string {
        h, _ := os.Hostname()
        return h
    })
    
    // Creating a "child" logger. In this example the child logger inherits the
    // settings from the parent logger, but has it's own name.
    logger = xlog.New("testing")
    child := logger.New("child")
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

The `xlog.GetLogger()` function gives you the ease of a global logger, with
the flexibility of named loggers. The function acts like a factory and global
registry, which is useful when you have hundreds or even thousands of objects
that need a logger, but you don't want to create thousands of *DefaultLogger instances.

```go
package main

import (
    "fmt"
    "github.com/dulo-tech/xlog"
)

func main() {
    // The first call to xlog.GetLogger() with the name "a" creates a
    // new logger for the name and returns it. The next time you call the
    // function with the same name, the same *DefaultLogger instance is returned.
    loggerA := xlog.GetLogger("a")
    loggerA.Append("stdout", xlog.DebugLevel)
    if loggerA == xlog.GetLogger("a") {
        fmt.Println("xlog.GetLogger() returned the same instance for the same name.")
    }
    
    loggerB := xlog.GetLogger("b")
    loggerB.Append("stdout", xlog.DebugLevel)
    if loggerB == xlog.GetLogger("b") {
        fmt.Println("These loggers have the same name, and same instance.")
    }
    
    if xlog.GetLogger("a") != xlog.GetLogger("b") {
        fmt.Println("These are different instances.")
    }
    
    // Outputs: 2014-11-15 09:40:28.693 a.DEBUG Test debug message.
    loggerA.Debug("Test debug message.")
    
    // Outputs: 2014-11-15 09:40:28.701 b.WARNING Test warning message.
    loggerB.Warning("Test warning message.")
    
    // Outputs: 2014-11-15 c.WARNING Test warning message.
    loggerC.Warning("Test warning message.")
    
    // This doesn't output anything because xlog.WarningLevel was set as the
    // global default level when loggerD was created.
    loggerD.Debug("Test debug message.")
}
```


#### Custom Formatters
You can create your own message formatter by creating a struct that implements
the `xlog.Formatter` interface, which has the following signature:

```go
type Formatter interface {
    // SetFormat changes the set message format.
    SetFormat(format string)
    
    // PlaceholderFunc adds a callback function which provides a replacement for key in a string format.
    PlaceholderFunc(key string, f func(string) string)
    
	// Format formats a log message for the given level.
	Format(name string, level Level, v ...interface{}) string
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
}

// SetFormat changes the set message format.
func (f *NullFormatter) SetFormat(format string) {
}

// PlaceholderFunc adds a callback function which provides a replacement for key in a string format.
func PlaceholderFunc(key string, f func(string) string) {
}

// Format formats a log message for the given level.
func (f *NullFormatter) Format(name string, level Level, v ...interface{}) string {
    return ""
}

func main() {
    // Creating a logger which discards all messages.
    logger := xlog.New("testing")
    logger.Append("stdout", xlog.DebugLevel)
    logger.Settings.Formatter := &NullFormatter{}
    
    // You can also assign the custom formatter to the global logger.
    xlog.SetFormatter(formatter)
}
```


#### Loggable Interface
The `xlog.New()` method and other New methods return an instance of
the struct `xlog.DefaultLogger`, which implements the `xlog.Loggable` interface.

```go
type Loggable interface {
	Writable() bool
	Closed() bool
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
	Writer(level Level) *LoggerWriter
}
```

Note the absence of methods like `Append()`, `MultiAppend()`, `AppendWriter()`,
and `Close()`, which are members of the `xlog.DefaultLogger` struct. The `xlog.Loggable`
interface does not concern itself for how a logger is configured or it's
lifecycle. The `xlog.Loggable` interface only exposes methods for logging
messages at various log levels.

You are encouraged to reference `xlog.Loggable` instead of `xlog.DefaultLogger` in your
code to keep it flexible to future API changes. For example instead of creating
a struct using `xlog.DefaultLogger` like this:

```go
type WebScraper struct {
    logger xlog.DefaultLogger
}

func NewWebScraper(logger *xlog.DefaultLogger) {
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
