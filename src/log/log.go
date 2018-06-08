package log

import (
	"time"
)

type ILogger interface {
	Init() error
	WriteMsg(when time.Time, msg string, level int) error
	Destory()
	Flush()
}
