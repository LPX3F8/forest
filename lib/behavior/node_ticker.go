package behavior

type BaseTicker struct {
	errs []error
}

func NewBaseTicker() *BaseTicker {
	return &BaseTicker{[]error{}}
}

func (*BaseTicker) TickerName() string {
	return "default"
}

func (*BaseTicker) OnBefore() (status Status, skip bool) {
	return Running, false
}

func (b *BaseTicker) OnAfter(status Status) Status {
	return status
}

func (b *BaseTicker) SetError(err ...error) {
	for _, e := range err {
		if e == nil {
			continue
		}
		b.errs = append(b.errs, e)
	}
}

func (b *BaseTicker) OnTick() Status {
	panic("implement me")
}

func (b *BaseTicker) Errors() []error {
	return b.errs
}
