package blackboard

import (
	"github.com/LPX3F8/froest/lib/store"
)

type Blackboard struct {
	namespace string
	scope     string
	store     store.Store
}

func NewBlackboard(namespace, scope string, store store.Store) *Blackboard {
	return &Blackboard{
		namespace: namespace,
		scope:     scope,
		store:     store,
	}
}

func (b *Blackboard) Set(k any, v any) error {
	return b.store.Set(k, v)
}

func (b *Blackboard) Get(k any) (any, bool) {
	return b.store.Get(k)
}

func GetValue[T any](bb *Blackboard, key any) (T, bool) {
	return store.GetValue[T](bb.store, key)
}
