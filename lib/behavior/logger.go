package behavior

import (
	"os"
	"path"
	"sync"

	"go.uber.org/zap"

	"github.com/LPX3F8/froest/log"
)

var (
	traceLogger   log.LevelLogger
	defaultLogger *zap.SugaredLogger
	mu            = new(sync.Mutex)
)

func init() {
	defaultLogger = log.NewLogger(
		path.Join(os.TempDir(), "trace_node.log"),
		log.Debug,
		true,
		log.WithWriteLogInFile(false),
		log.WithCaller(false),
	)
	traceLogger = defaultLogger
}

func SetTraceLogger(logger log.LevelLogger) {
	mu.Lock()
	defer mu.Unlock()
	traceLogger = logger
}

func ResetTraceLogger() {
	mu.Lock()
	defer mu.Unlock()
	traceLogger = defaultLogger
}
