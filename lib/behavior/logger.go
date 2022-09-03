package behavior

import (
	"go.uber.org/zap"

	"github.com/LPX3F8/froest/log"
)

var traceLogger log.LevelLogger

func init() {
	traceLogger = log.NewLogger("", zap.DebugLevel, true)
}
