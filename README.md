xLog
====
Simple logger for the Go language.

Documentation is available from the [GoDoc website](http://godoc.org/github.com/dulo-tech/xlog).


#### Installation
`go get github.com/dulo-tech/xlog`


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
    // name will appear in the logged messages.
    logger := xlog.NewLogger("testing")
    
    // Add files where the logs will be written. You can use a full system path
    // to a file, or the aliases "stdout", "stderr", and "stdin". Each file you
    // add is assigned a logging level. The file is written to when a message is
    // logged to that level or a greater level. The levels are:
    //
    // xlog.Debug
    // xlog.Info
    // xlog.Notice
    // xlog.Warning
    // xlog.Error
    // xlog.Critical
    // xlog.Alert
    // xlog.Emergency
    //
    // When you append a file using the xlog.Debug level, the file will be
    // written to for every level, because every level is greater than xdt.Debug.
    logger.Append("stdout", xlog.Debug)
    
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
    
    // Any log message with the xlog.Warning level and above will be logged to
    // stdout.
    logger = xlog.NewLogger("testing")
    logger.Append("stdout", xlog.Warning)
    
    // This doesn't output anything because the xlog.Debug level is lower than
    // the xlog.Warning level.
    logger.Debug("Test debug message.")
    
    // Outputs: 2014-11-15 09:40:28.701 testing.WARNING Test warning message.
    logger.Warning("Test warning message.")
    
    // This doesn't output anything because the Info level is lower than the
    // Warning level.
    logger.Infof("Test %s message.", "info")
    
    // You can append as many files as you want. This logs messages xlog.Debug and
    // above to stdout, and messages xlog.Error level and above to a file.
    // Note, when you have the logger open files, the files need to be closed by
    // the logger by deferring the defer logger.Close() method.
    logger = xlog.NewLogger("testing")
    logger.Append("stdout", xlog.Debug)
    logger.Append("/var/logs/main-error.log", xlog.Error)
    defer logger.Close()
    
    // You can manage the files yourself by using the logger.AppendWriter()
    // method.
    fp, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE | os.O_APPEND, 0666)
    if err != nil {
        panic(err)
    }
    defer fp.Close()
    logger.AppendWriter(fp)
    
    // You can replicate the functionality of Go's system logger log.Fatal()
    // and log.Panic() using logger.FatalOn and logger.PanicOn.
    
    // Logging a message to the xlog.Critical level will cause a fatal shut down
    // using os.Exit(1).
    logger.FatalOn = xlog.Critical
    
    // Logging a message to either xlog.Alert or xlog.Emergency levels causes a
    // panic using panic().
    logger.PanicOn = xlog.Alert | xlog.Emergency
    
    // Change the way the log messages are formatted. The xlog.Formatter interface
    // requires a format string and a name. The format string defines how the
    // log messages are formatted.
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
        "{date} {name}.{level} - {message}",
        "testing",
    )
    logger = xlog.NewFormattedLogger(formatter)
    logger.Append("stdout", xlog.Debug)
    
    // Outputs: 2014-11-15 09:59:32.427 testing.DEBUG - Test debug message.
    logger.Debug("Test debug message.")
    
    // The message format can be changed without setting a new Formatter. The
    // name can be changed as well.
    logger.Formatter.SetMessageFormat("{date} {message}")
    logger.Formatter.SetName("debug-testing")
}
```
