package log

import (
	"fmt"
	"os"
)

//Console logger, print log message to console, and each message level has
//different color
type console struct {
	level  int
	colors []string
	cache  chan *Record
	quit   chan bool
}

//Create new console logger
func NewConsoleLogger(level int) Logger {
	c := &console{
		level: level,
		colors: []string{
			DEBUG: "\033[32m",
			INFO:  "\033[36m",
			WARN:  "\033[33m",
			ERROR: "\033[31m",
			FATAL: "\033[35m",
			PANIC: "\033[37m",
		},
		cache: make(chan *Record, BUFFER_CAPACITY),
		quit:  make(chan bool),
	}

	go c.run()

	return c
}

//run is goroutine which actually handle the log message
func (c *console) run() {

	for {
		select {
		case rec := <-c.cache:
			c.write(rec)
		case <-c.quit:
			return
		}
	}
}

//flush print all the buffer message to console
func (c *console) flush() {

	for {
		select {
		case rec := <-c.cache:
			c.write(rec)
		default:
			return
		}
	}
}

func (c *console) write(record *Record) {

	timestr := record.timestamp.Format("2006/01/02 15:04:04")

	fmt.Fprintln(os.Stdout, c.colors[record.level], timestr, "[", levelstr[record.level], "]:", record.message, "\033[0m")
}

func (c *console) Level() int {
	return c.level
}

func (c *console) Write(record *Record) {
	c.cache <- record
}

func (c *console) Close() {
	c.quit <- true
	c.flush()
	close(c.cache)
	close(c.quit)
}
