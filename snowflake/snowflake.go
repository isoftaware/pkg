package snowflake

import (
	"time"

	"github.com/sony/sonyflake"
)

const (
	timeLayout = "2006-01-02 15:04:05"
	startTime  = "2020-10-01 00:00:00"
)

var sf *sonyflake.Sonyflake

func init() {
	st, _ := time.Parse(timeLayout, startTime)

	sf = sonyflake.NewSonyflake(sonyflake.Settings{
		StartTime: st,
	})
}

func NextID() (uint64, error) {
	return sf.NextID()
}
