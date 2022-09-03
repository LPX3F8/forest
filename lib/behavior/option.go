package behavior

type TreeOption struct {
	EnableTrace bool
}

func (o *TreeOption) debugf(template string, arg ...interface{}) {
	if !o.EnableTrace {
		return
	}
	traceLogger.Debugf(template, arg...)
}

func (o *TreeOption) debug(template string, arg ...interface{}) {
	if !o.EnableTrace {
		return
	}
	traceLogger.Debugln(arg...)
}
