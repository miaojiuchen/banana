package log

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// FileLogWriter todo
// type FileLogWriter struct {
// 	fileLogWriter
// }

type fileLogWriter struct {
	sync.RWMutex

	filename   string
	fileWriter *os.File

	maxLines int
	curLine  int

	maxSize int
	curSize int

	daily    bool
	openTime time.Time

	level int
}

func newFileLogWriter() ILogger {
	return &fileLogWriter{}
}

func (w *fileLogWriter) Init() error {
	fmt.Println("init filelogwriter")

	if len(w.filename) == 0 {
		return errors.New("must specify a filename")
	}

	return w.startLogger()
}

func (w *fileLogWriter) startLogger() error {
	file, err := w.ensureLogFile()
	if err != nil {
		return err
	}

	if w.fileWriter != nil {
		w.fileWriter.Close()
	}

	w.fileWriter = file

	return w.initFd()
}

func (w *fileLogWriter) ensureLogFile() (*os.File, error) {
	perm := 0644

	fd, err := os.OpenFile(w.filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(perm))
	if err == nil {
		os.Chmod(w.filename, os.FileMode(perm))
	}
	return fd, err
}

func (w *fileLogWriter) initFd() error {
	fd := w.fileWriter

	fInfo, err := fd.Stat()
	if err != nil {
		return fmt.Errorf("get stat err: %s", err)
	}

	w.curSize = int(fInfo.Size())
	w.curLine = 0

	if fInfo.Size() > 0 && w.maxLines > 0 {
		count, err := w.lines()
		if err != nil {
			return err
		}
		w.curLine = count
	}

	return nil
}

func (w *fileLogWriter) lines() (int, error) {
	fd, err := os.Open(w.filename)
	if err != nil {
		return 0, err
	}
	defer fd.Close()

	buf := make([]byte, 32768) // 32k
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := fd.Read(buf)
		if err != nil && err != io.EOF {
			return count, err
		}

		count += bytes.Count(buf[:c], lineSep)

		if err == io.EOF {
			break
		}
	}

	return count, nil
}

func (w *fileLogWriter) WriteMsg(when time.Time, msg string, level int) error {
	if level > w.level {
		return nil
	}

	timeHeader := formatTimeHeader(when)

	msg = timeHeader + msg + "\n"

	w.RLock()

	if w.needRotate(len(msg)) {
		w.RUnlock()
		w.Lock()

		// check onemore time in case of modify
		if w.needRotate(len(msg)) {
			err := w.doRotate(when)
			if err != nil {
				fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", w.filename, err)
			}
		}

		w.Unlock()
	} else {
		w.RUnlock()
	}

	w.Lock()

	_, err := w.fileWriter.Write([]byte(msg))
	if err == nil {
		w.curLine++
		w.curSize += len(msg)
	}

	w.Unlock()

	return err
}

func (w *fileLogWriter) needRotate(msgSize int) bool {
	return (w.maxLines > 0 && w.curLine >= w.maxLines) ||
		(w.maxSize > 0 && w.curSize >= w.maxSize)
}

func (w *fileLogWriter) doRotate(when time.Time) error {
	return nil
}

func (w *fileLogWriter) Destory() {

}

func (w *fileLogWriter) Flush() {

}
