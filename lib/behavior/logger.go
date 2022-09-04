package behavior

import (
	"os"
	"path"
	"sync"

	"github.com/LPX3F8/froest/log"
)

var (
	traceLogger log.LevelLogger
	mu          = new(sync.Mutex)
)

func init() {
	traceLogger = log.NewLogger(
		path.Join(os.TempDir(), "trace_node.log"),
		log.Debug,
		true,
		log.WithWriteLogInFile(false),
		log.WithCaller(false),
	)
}

func SetTraceLogger(logger log.LevelLogger) {
	mu.Lock()
	defer mu.Unlock()
	traceLogger = logger
}
