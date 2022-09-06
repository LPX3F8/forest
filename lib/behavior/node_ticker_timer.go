package behavior

import (
	"time"

	"github.com/LPX3F8/orderedmap"
)

type Period struct {
	Label     string
	StartTime time.Time     `json:"startTime"`
	EndTime   time.Time     `json:"endTime"`
	Duration  time.Duration `json:"duration"`
}

func (p *Period) Start() *Period {
	p.StartTime = time.Now()
	return p
}

func (p *Period) Stop() *Period {
	p.EndTime = time.Now()
	p.Duration = p.EndTime.Sub(p.StartTime)
	return p
}

type Timer struct {
	node    IBTreeNode
	periods *orderedmap.OrderedMap[string, *Period]
}

func NewSimpleTimer(node IBTreeNode) ITimer {
	return &Timer{
		node:    node,
		periods: orderedmap.New[string, *Period](),
	}
}

func (*Timer) TimerName() string {
	return "default"
}

func (t *Timer) Start(label string) {
	p := Period{Label: label}
	if s, ok := t.periods.Load(label); ok {
		s.Start()
		return
	}
	t.periods.Store(label, p.Start())
}

func (t *Timer) Stop(label string) {
	if s, ok := t.periods.Load(label); ok {
		s.Stop()
	}
}

func (t *Timer) Time(label string, f func()) {
	t.Start(label)
	defer t.Stop(label)
	f()
}

func (t *Timer) Report() []*Period {
	return t.periods.Slice()
}
