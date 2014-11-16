xLog
====
Simple logger for the Go language.

Documentation is available from the [GoDoc website](http://godoc.org/github.com/dulo-tech/xlog).


#### Installation
`go get github.com/dulo-tech/xlog`


#### Examples
```go
package main

import "github.com/dulo-tech/xlog"

func main() {
    // Any log message with the Debug level and above will be logged to stdout.
    logger := xlog.NewLogger()
    logger.Append("stdout", xlog.Debug)
    
    // Outputs: 2014-11-15 09:40:28.693 [DEBUG] Test debug message.
    logger.Debug("Test debug message.")
    
    // Outputs: 2014-11-15 09:40:28.701 [WARNING] Test warning message.
    logger.Warning("Test warning message.")
    
    // Outputs: 2014-11-15 09:40:28.723 [INFO] Test info message.
    logger.Infof("Test %s message.", "info")
    
    // Any log message with the Warning level and above will be logged to stdout.
    logger = xlog.NewLogger()
    logger.Append("stdout", xlog.Warning)
    
    // This doesn't output anything because the Debug level is lower than the
    // Warning level.
    logger.Debug("Test debug message.")
    
    // Outputs: 2014-11-15 09:40:28.701 [WARNING] Test warning message.
    logger.Warning("Test warning message.")
    
    // This doesn't output anything because the Info level is lower than the
    // Warning level.
    logger.Infof("Test %s message.", "info")
    
    // This logs messages Debug and above to stdout, and messages Error level and
    // above to a file.
    logger = xlog.NewLogger()
    logger.Append("stdout", xlog.Debug)
    logger.Append("/var/logs/main-error.log", Logs.Error)
    
    // Logging a message to the Critical level will cause a fatal shut down
    // using os.Exit(1).
    logger.FatalOn = xlog.Critical
    
    // Logging a message to either Alert or Emergency levels causes a panic
    // using panic().
    logger.PanicOn = xlog.Alert | xlog.Emergency
    
    // Change the way the log messages are formatted.
    logger.Formatter = xlog.NewDefaultFormatter("{date} main.{level} - {message}")
    
    // Outputs: 2014-11-15 09:54:16.278 main.DEBUG - This is a debug test.
    logger.Debug("Test debug message.")
    
    // Format the output date using Go's date format. 
    // See: http://golang.org/pkg/time/#Time.Format
    logger.Formatter = xlog.NewDefaultFormatter("{date|Jan _2 15:04:05} {level} {message}")
    
    // Outputs: Nov 15 09:56:56 DEBUG Test debug message.
    logger.Debug("Test debug message.")
    
    // Creating a logger with a pre-configured formatter.
    formatter := xlog.NewDefaultFormatter("{date} main.{level} - {message}")
    logger = xlog.NewFormattedLogger(formatter)
    logger.Append("stdout", xlog.Debug)
    
    // Outputs: 2014-11-15 09:59:32.427 main.DEBUG - Test debug message.
    logger.Debug("Test debug message.")
}
```
