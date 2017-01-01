package log

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type file struct {
	level    int
	path     string
	duration time.Duration
	filename string
	file     *os.File
	cache    chan *Record
	quit     chan bool
}

func NewFileLogger(level int, args ...interface{}) Logger {

	path := ""

	if len(args) == 1 {
		ok := false
		if path, ok = args[0].(string); !ok {
			path = ""
		}
	}

	os.MkdirAll(path, 0770)

	f := &file{
		level:    level,
		path:     path,
		filename: "",
		file:     nil,
		cache:    make(chan *Record, BUFFER_CAPACITY),
		quit:     make(chan bool),
	}

	go f.run()

	return f
}

func (f *file) run() {

	f.rotate()

	for {
		select {
		case rec := <-f.cache:
			f.write(rec)
		case <-time.After(f.duration):
			f.rotate()
			go f.sweep()
		case <-f.quit:
			return
		}
	}
}

func (f *file) rotate() error {

	if f.file != nil {
		f.file.Close()
	}

	f.filename = time.Now().Format("2006-01-02") + ".log"

	var err error = nil
	fname := f.path + "/" + f.filename
	f.file, err = os.OpenFile(fname, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		return err
	}

	day := time.Now().AddDate(0, 0, 1)
	day = time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.Local)
	f.duration = day.Sub(time.Now())

	return nil
}

func (f *file) sweep() {

	filepath.Walk(f.path, func(path string, fi os.FileInfo, err error) error {
		if fi == nil {
			return err
		}

		if fi.ModTime().Add(time.Hour * 24 * 30).Before(time.Now()) {
			os.Remove(f.path + "/" + path)
		}

		return nil
	})
}

func (f *file) write(record *Record) {

	timestr := record.timestamp.Format("2006/01/02 15:04:04")

	fmt.Fprintln(f.file, timestr, "[", levelstr[record.level], "]:", record.message)
}

func (f *file) flush() {
	for {
		select {
		case rec := <-f.cache:
			f.write(rec)
		default:
			return
		}
	}
}

func (f *file) Level() int {
	return f.level
}

func (f *file) Write(record *Record) {
	f.cache <- record
}

func (f *file) Close() {
	defer func() {
		if f.file != nil {
			f.file.Close()
		}
	}()

	f.quit <- true
	f.flush()
	close(f.cache)
	close(f.quit)
}
