package log

import "go.uber.org/zap/zapcore"

const (
	defLogFileName     = "forest.log"
	defMaxSize         = 512
	defMAxBackups      = 2
	defMaxAges         = 7
	defEnableCaller    = true
	defEnableStdout    = true
	defEnableCompress  = false
	defEnableWriteFile = false
	defLogLevel        = zapcore.DebugLevel
)
