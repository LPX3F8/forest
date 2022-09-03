package log

import (
	"os"
	"path"
	"sync"

	"github.com/natefinch/lumberjack"
	"go.uber.org/atomic"
	"go.uber.org/zap/zapcore"
)

type ZapConfig struct {
	logFileName  atomic.String // logfile path, default os.TempDir()
	logLevel     atomic.Value
	maxSize      atomic.Int64 // logfile max size in MB, default 100MB
	maxBackups   atomic.Int64 // logfile backups, default unlimited
	maxAges      atomic.Int64 // logfile max save time of days, default unlimited
	withCompress atomic.Bool  // logfile backups compress, default false
	withStdout   atomic.Bool  // write log with console, default true
	withCaller   atomic.Bool  // print log with caller, default true
	*sync.Mutex
}

// GetLogFilePath return the path of logfile
func (c ZapConfig) GetLogFilePath() string {
	return path.Join(c.logFileName.Load())
}

// GetLogLevel return the log level
func (c ZapConfig) GetLogLevel() zapcore.Level {
	return c.logLevel.Load().(zapcore.Level)
}

// WithCaller return need if print caller
func (c ZapConfig) WithCaller() bool {
	return c.withCaller.Load()
}

// WithStdout return need if print into os.stdout
func (c ZapConfig) WithStdout() bool {
	return c.withStdout.Load()
}

// GetLogFileRoller return the pointer of lumberjack.Logger
func (c ZapConfig) GetLogFileRoller() *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   c.GetLogFilePath(),
		MaxSize:    int(c.maxSize.Load()),
		MaxBackups: int(c.maxBackups.Load()),
		MaxAge:     int(c.maxAges.Load()),
		Compress:   c.withCompress.Load(),
	}
}

// GetZapWriteSyncer return the writeSyncer with config
func (c ZapConfig) GetZapWriteSyncer() zapcore.WriteSyncer {
	ws := zapcore.AddSync(c.GetLogFileRoller())
	if c.WithStdout() {
		ws = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), ws)
	}
	return ws
}

func NewZapLogConfig(options ...Option) *ZapConfig {
	c := &ZapConfig{Mutex: new(sync.Mutex)}

	// default options
	options = append([]Option{
		WithLogFile(path.Join(os.TempDir(), defaultLogFile)),
		WithMaxSize(512),
		WithMaxBackups(2),
		WithMaxAges(7),
		WithCaller(true),
		WithStdout(true),
		WithLevel(zapcore.DebugLevel),
	}, options...)

	for _, o := range options {
		o(c)
	}
	return c
}
