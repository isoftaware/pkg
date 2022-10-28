package log

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type logger struct {
	option  *Option
	pool    *sync.Pool
	buffer  chan *log
	version int64
	output  *output
}

func newLogger(option *Option) (*logger, error) {
	var err error

	lgr := &logger{
		option: option,
		pool: &sync.Pool{
			New: func() any {
				return newLog()
			},
		},
		buffer: make(chan *log, maxBufferSize),
	}

	t := time.Now()
	lgr.version = t.Unix() / int64(lgr.option.Shard)
	lgr.output, err = newOutput(option.OutputDir, lgr.genKey(t))
	if err != nil {
		return nil, err
	}

	go lgr.consumer()

	return lgr, nil
}

func (lgr *logger) consumer() {
	defer func() {
		_ = lgr.output.Flush()

		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "[micro-log] logger.consumer recovered from panic, err: %v\n", err)
			go lgr.consumer()
		}
	}()

	for {
		select {
		case lg, ok := <-lgr.buffer:
			if !ok {
				return
			}
			lgr.printLog(lg)
			lgr.putLog(lg)
		case <-time.Tick(time.Second):
			err := lgr.output.Flush()
			if err != nil {
				fmt.Fprintf(os.Stderr, "[micro-log] failed to flush output buffer, err: %v\n", err)
			}
		}
	}
}

func (lgr *logger) getLog() *log {
	return lgr.pool.Get().(*log)
}

func (lgr *logger) putLog(lg *log) {
	lg.Reset()
	lgr.pool.Put(lg)
}

func (lgr *logger) appendBuffer(lg *log) {
	select {
	case lgr.buffer <- lg:
	default:
	}
}

func (lgr *logger) printLog(lg *log) {
	version := lg.Time.Unix() / int64(lgr.option.Shard)
	if version > lgr.version {
		key := lgr.genKey(lg.Time)
		err := lgr.output.ResetFile(key)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[micro-log] failed to reset output, key=%s, err: %v\n", key, err)
		} else {
			lgr.version = version
		}
	}

	err := lgr.output.Writeln(lg.Dump())
	if err != nil {
		fmt.Fprintf(os.Stderr, "[micro-log] failed to print log to output, err: %v\n", err)
	}
}

func (lgr *logger) genKey(t time.Time) string {
	switch lgr.option.Shard {
	case ShardPerHour:
		return t.Format(timeFormatHour)
	case ShardPerDay:
		return t.Format(timeFormatDay)
	}
	return t.Format(timeFormatDay)
}

func (lgr *logger) Error(err error, message string) {
	if lgr.option.EnabledLevels&LevelError == 0 {
		return
	}

	lg := lgr.getLog()
	lg.Fill(LevelError, err, message)

	lgr.appendBuffer(lg)
}

func (lgr *logger) Warn(message string) {
	if lgr.option.EnabledLevels&LevelWarn == 0 {
		return
	}

	lg := lgr.getLog()
	lg.Fill(LevelWarn, nil, message)

	lgr.appendBuffer(lg)
}

func (lgr *logger) Info(message string) {
	if lgr.option.EnabledLevels&LevelInfo == 0 {
		return
	}

	lg := lgr.getLog()
	lg.Fill(LevelInfo, nil, message)

	lgr.appendBuffer(lg)
}

func (lgr *logger) Debug(message string) {
	if lgr.option.EnabledLevels&LevelDebug == 0 {
		return
	}

	lg := lgr.getLog()
	lg.Fill(LevelDebug, nil, message)

	lgr.appendBuffer(lg)
}
