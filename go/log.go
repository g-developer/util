package util

import (
	"runtime"
	"time"
	"sync"
	"os"
	"fmt"
)

type log struct {
	level int
	filename string
	fd *os.File
	maxSize int
	shard int
	suffix string
}

var logOnce sync.Once
var logIns *log = nil

const (
	FATAL = 1
	WARN = 2
	NOTICE = 4
	DEBUG = 8
	INFO = 16
)

const (
	LOG_SHARD_BY_HOUR = 1
	LOG_SHARD_BY_DAY = 2
	LOG_NO_SHARD = 4
)

func newLog () *log{
	return &log{INFO, "", nil, -1, LOG_NO_SHARD, ".log"}
}

func GetLogInstance () *log {
	logOnce.Do(func() {
		logIns = newLog()
	})
	return logIns
}

func (self *log) SetFile (filename string) error {
	self.filename = filename + self.suffix
	f, err := os.OpenFile(self.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if nil == err {
		self.fd = f
	}
	return err
}

func (self *log) SetLevel (level int) {
	if level >= INFO {
		self.level = INFO
	} else if level < INFO && level >= DEBUG {
		self.level = DEBUG
	} else if level < DEBUG && level >= NOTICE {
		self.level = NOTICE
	} else if level < NOTICE && level >= WARN {
		self.level = WARN
	} else {
		self.level = FATAL
	}
}


func (self *log) writeLog (level string, format string, v ... interface{}) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	//_, filename := path.Split(file)
	tmp := time.Now().Format("2006-01-02 15:04:05.000") + " " + fmt.Sprintf("[%-6s] [%s:%d] msg : ", level, file, line) + fmt.Sprintf(format, v...)
	fmt.Println(tmp)
}

func (self *log) Fatal (format string, v ...interface{}) {
	if FATAL > self.level {
		return
	}
	self.writeLog("FATAL", format, v...)
}

func (self *log) Warn (format string, v ...interface{}) {
	if WARN > self.level {
		return
	}
	self.writeLog("WARN", format, v...)
}

func (self *log) Notice (format string, v ...interface{}) {
	if NOTICE > self.level {
		return
	}
	self.writeLog("NOTICE", format, v...)
}

func (self *log) Debug (format string, v ...interface{}) {
	if DEBUG > self.level {
		return
	}
	self.writeLog("DEBUG", format, v...)
}

func (self *log) Info (format string, v ...interface{}) {
	if INFO > self.level {
		return
	}
	self.writeLog("INFO", format, v...)
}
