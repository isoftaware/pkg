package log

import (
	"reflect"
	"runtime"
	"strconv"
	"time"
	"unsafe"

	"github.com/bytedance/sonic"
)

type log struct {
	Level   string    `json:"level"`
	Error   string    `json:"error,omitempty"`
	Message string    `json:"message"`
	Stack   []string  `json:"stack,omitempty"`
	Time    time.Time `json:"time"`
}

func newLog() *log {
	return &log{}
}

func (lg *log) Fill(level level, err error, message string) {
	switch level {
	case LevelError:
		lg.Level = levelErrorName
	case LevelWarn:
		lg.Level = levelWarnName
	case LevelInfo:
		lg.Level = levelInfoName
	case LevelDebug:
		lg.Level = levelDebugName
	}

	lg.Time = time.Now()
	lg.Message = message

	if err != nil {
		lg.Error = err.Error()
	}

	if level == LevelError {
		lg.dumpStack()
	}
}

func (lg *log) dumpStack() {
	if cap(lg.Stack) == 0 {
		lg.Stack = make([]string, 0, maxStackDepth)
	}

	pc := make([]uintptr, maxStackDepth)
	runtime.Callers(baseCallersSkip, pc)
	frames := runtime.CallersFrames(pc)
	for {
		frame, more := frames.Next()
		lg.Stack = append(lg.Stack, frame.File+":"+strconv.Itoa(frame.Line))
		if !more {
			break
		}
	}
}

func (lg *log) Dump() []byte {
	dump, _ := sonic.Marshal(lg)
	return dump
}

func (lg *log) Reset() {
	lg.Error = ""
	if len(lg.Stack) != 0 {
		(*reflect.SliceHeader)(unsafe.Pointer(&lg.Stack)).Len = 0
	}
}
