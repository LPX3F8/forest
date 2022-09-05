package behavior

import (
	"sync"

	"github.com/LPX3F8/froest/lib/store"
)

var manager *bbManager

func init() {
	manager = &bbManager{
		bb:      map[string]map[string]*Blackboard{},
		RWMutex: new(sync.RWMutex),
	}
}

type bbManager struct {
	bb map[string]map[string]*Blackboard
	*sync.RWMutex
}

func (m *bbManager) GetBlackboard(namespace, scope string, s store.Store) *Blackboard {
	m.Lock()
	defer m.Unlock()
	if bb, ok := m.bb[namespace]; ok {
		if cc, ok := bb[scope]; ok {
			return cc
		}
	}
	m.bb[namespace] = map[string]*Blackboard{}
	m.bb[namespace][scope] = NewBlackboard(namespace, scope, s)
	return m.bb[namespace][scope]
}

func (m *bbManager) GetDefaultBlackboard(namespace, scope string) *Blackboard {
	return manager.GetBlackboard(namespace, scope, store.NewMemStore())
}

func TreeBlackboard(namespace, scope string) *Blackboard {
	return manager.GetDefaultBlackboard(namespace, scope)
}
