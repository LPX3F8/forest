package behavior

import (
	"sync"

	"github.com/LPX3F8/froest/log"
)

var (
	traceLogger log.LevelLogger
	mu          = new(sync.Mutex)
)

func init() {
	traceLogger = log.NewLogger("", log.Debug, true)
}

func SetTraceLogger(logger log.LevelLogger) {
	mu.Lock()
	defer mu.Unlock()
	traceLogger = logger
}
