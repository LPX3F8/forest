package log

import (
	"os"
	"sync"

	"github.com/natefinch/lumberjack"
	"go.uber.org/atomic"
	"go.uber.org/zap/zapcore"
)

type ZapConfig struct {
	logFileName  *atomic.String // logfile path, default os.TempDir()
	logLevel     *atomic.String
	maxSize      *atomic.Int64 // logfile max size in MB, default 512MB
	maxBackups   *atomic.Int64 // logfile backups, default unlimited
	maxAges      *atomic.Int64 // logfile max save time of days, default unlimited
	withCompress *atomic.Bool  // logfile backups compress, default false
	withStdout   *atomic.Bool  // write log with console, default true
	withFile     *atomic.Bool  // write log into file, default false
	withCaller   *atomic.Bool  // print log with caller, default true
	*sync.Mutex
}

func NewDefaultZapConfig() *ZapConfig {
	return &ZapConfig{
		logFileName:  atomic.NewString(defLogFileName),
		logLevel:     atomic.NewString(Debug.String()),
		maxSize:      atomic.NewInt64(defMaxSize),
		maxBackups:   atomic.NewInt64(defMAxBackups),
		maxAges:      atomic.NewInt64(defMaxAges),
		withCompress: atomic.NewBool(defEnableCompress),
		withStdout:   atomic.NewBool(defEnableStdout),
		withFile:     atomic.NewBool(defEnableWriteFile),
		withCaller:   atomic.NewBool(defEnableCaller),
		Mutex:        new(sync.Mutex),
	}
}

// GetLogFilePath return the path of logfile
func (c ZapConfig) GetLogFilePath() string {
	return c.logFileName.Load()
}

// GetLogLevel return the log level
func (c ZapConfig) GetLogLevel() zapcore.Level {
	switch c.logLevel.Load() {
	case Debug.String():
		return zapcore.DebugLevel
	case Info.String():
		return zapcore.InfoLevel
	case Warn.String():
		return zapcore.WarnLevel
	case Error.String():
		return zapcore.ErrorLevel
	default:
		return defLogLevel
	}
}

// WithCaller return need if print caller
func (c ZapConfig) WithCaller() bool {
	return c.withCaller.Load()
}

// WithStdout return need if print into os.stdout
func (c ZapConfig) WithStdout() bool {
	return c.withStdout.Load()
}

// WithFile return need if write log into file
func (c ZapConfig) WithFile() bool {
	return c.withFile.Load()
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
	ws := make([]zapcore.WriteSyncer, 0)
	if c.WithFile() {
		ws = append(ws, zapcore.AddSync(c.GetLogFileRoller()))
	}
	if c.WithStdout() {
		ws = append(ws, zapcore.AddSync(os.Stdout))
	}

	return zapcore.NewMultiWriteSyncer(ws...)
}

func NewZapLogConfig(options ...Option) *ZapConfig {
	c := NewDefaultZapConfig()
	for _, o := range options {
		o(c)
	}
	return c
}
