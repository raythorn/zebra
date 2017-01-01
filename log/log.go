package log

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"time"
)

const (
	DEBUG = iota
	INFO
	WARN
	ERROR
	FATAL
	PANIC
)

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

type Record struct {
	level     int
	verbose   string
	message   string
	timestamp time.Time
}

type Logger interface {
	Level() int
	Write(record *Record)
	Close()
}

type logger struct {
	verbose bool
	channel map[string]Logger
}

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

func Close() {
	log4f.close()
}

func Debug(format string, args ...interface{}) {
	log4f.log(DEBUG, format, args...)
}

func Info(format string, args ...interface{}) {
	log4f.log(INFO, format, args...)
}

func Warning(format string, args ...interface{}) {
	log4f.log(WARN, format, args...)
}

func Error(format string, args ...interface{}) {
	log4f.log(ERROR, format, args...)
}

func Fatal(format string, args ...interface{}) {
	log4f.log(FATAL, format, args...)
	os.Exit(0)
}

func Panic(format string, args ...interface{}) {
	log4f.log(PANIC, format, args...)
	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}
	panic(msg)
}
