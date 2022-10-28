package log

type Option struct {
	OutputDir     string
	Shard         shard
	EnabledLevels level
}

func defaultOption() *Option {
	return &Option{
		OutputDir:     "logs",
		Shard:         ShardPerDay,
		EnabledLevels: LevelError | LevelWarn | LevelInfo,
	}
}
