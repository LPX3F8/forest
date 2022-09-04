package log

type Option func(*ZapConfig)

func WithLogFile(name string) Option {
	return func(config *ZapConfig) {
		config.logFileName.Store(name)
	}
}

func WithMaxSize(size int) Option {
	return func(config *ZapConfig) {
		config.maxSize.Store(int64(size))
	}
}

func WithMaxBackups(backups int) Option {
	return func(config *ZapConfig) {
		config.maxBackups.Store(int64(backups))
	}
}

func WithMaxAges(ages int) Option {
	return func(config *ZapConfig) {
		config.maxAges.Store(int64(ages))
	}
}

func WithCompress(compress bool) Option {
	return func(config *ZapConfig) {
		config.withCompress.Store(compress)
	}
}

func WithStdout(stdout bool) Option {
	return func(config *ZapConfig) {
		config.withStdout.Store(stdout)
	}
}

func WithLevel(level Level) Option {
	return func(config *ZapConfig) {
		config.logLevel.Store(level.String())
	}
}

func WithCaller(with bool) Option {
	return func(config *ZapConfig) {
		config.withCaller.Store(with)
	}
}
