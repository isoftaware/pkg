package log

type level uint

const (
	LevelError level = 1 << iota
	LevelWarn
	LevelInfo
	LevelDebug
)

const (
	levelErrorName = "ERROR"
	levelWarnName  = "WARN"
	levelInfoName  = "INFO"
	levelDebugName = "DEBUG"
)

type shard int64

const (
	ShardPerHour shard = 3600      // per hour
	ShardPerDay  shard = 3600 * 24 // per day
)

const (
	timeFormatHour = "2006-01-02_15"
	timeFormatDay  = "2006-01-02"
)

const (
	maxBufferSize   = 1 << 10
	baseCallersSkip = 5
	maxStackDepth   = 16
)
