package log

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"time"
)

//Log level
const (
	DEBUG = iota
	INFO
	WARN
	ERROR
	FATAL
	PANIC
)

//Capacity of buffer channel
const (
	BUFFER_CAPACITY = 16
)

var (
	log4f    *logger = nil
	levelstr         = [...]string{"Debug", "Info", "Warning", "Error", "Fatal", "Panic"}
)

func init() {
	log4f = &logger{verbose: false, channel: make(map[string]Logger)}
	Add("console", DEBUG)
}

//Log message record, minimal log dispatch unit
type Record struct {
	level     int
	verbose   string
	message   string
	timestamp time.Time
}

//Logger interface, all supported logger MUST implement this interface
type Logger interface {
	//Return logger's level
	Level() int
	//Write a message to logger
	Write(record *Record)
	//Close logger
	Close()
}

type logger struct {
	verbose bool
	channel map[string]Logger
}

//Internal log function
func (l *logger) log(level int, format string, args ...interface{}) {

	skip := true
	for _, logger := range l.channel {
		if logger.Level() <= level {
			skip = false
			break
		}
	}

	if skip {
		return
	}

	verbose := ""
	if l.verbose {
		pc, file, fileno, ok := runtime.Caller(2)
		if ok {
			verbose = fmt.Sprintf("%s:%d - %s", file, fileno, runtime.FuncForPC(pc).Name())
		}
	}

	message := format
	if len(args) > 0 {
		message = fmt.Sprintf(format, args...)
	}

	record := &Record{
		level:     level,
		verbose:   verbose,
		message:   message,
		timestamp: time.Now(),
	}

	for _, logger := range l.channel {

		if logger.Level() <= level {
			logger.Write(record)
		}
	}
}

func (l *logger) close() {

	for key, logger := range l.channel {
		logger.Close()
		delete(l.channel, key)
	}
}

//Log package default add one logger(console), means default all the log message will
//print to console. But you can use this function to add new logger to log engine, and
//current only two logger supported, "console" and "file".
func Add(logger string, level int, args ...interface{}) error {

	var err error = nil

	switch logger {
	case "console":
		c := NewConsoleLogger(level)
		log4f.channel[logger] = c
	case "file":
		f := NewFileLogger(level, args...)
		log4f.channel[logger] = f
	default:
		err = errors.New("Logger not supported!")
	}

	return err
}

//Close log engine
func Close() {
	log4f.close()
}

//Debug print debug message
func Debug(format string, args ...interface{}) {
	log4f.log(DEBUG, format, args...)
}

//Info print infomation message
func Info(format string, args ...interface{}) {
	log4f.log(INFO, format, args...)
}

//Warning print warning message
func Warning(format string, args ...interface{}) {
	log4f.log(WARN, format, args...)
}

//Error print error message
func Error(format string, args ...interface{}) {
	log4f.log(ERROR, format, args...)
}

//Fatal print fatal error message, and app will quit if this function called
func Fatal(format string, args ...interface{}) {
	log4f.log(FATAL, format, args...)
	os.Exit(0)
}

//Panic print panic message, and app will trigger panic message if called
func Panic(format string, args ...interface{}) {
	log4f.log(PANIC, format, args...)
	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}
	panic(msg)
}
