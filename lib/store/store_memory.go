package store

import (
	"errors"
	"fmt"
	"sync"

	j "github.com/json-iterator/go"
)

var json = j.ConfigCompatibleWithStandardLibrary

const InMemStore = "IN_MEMORY_STORE"

type Memory struct {
	store *sync.Map
	*sync.RWMutex
}

func NewMemStore() Store {
	return &Memory{
		store:   new(sync.Map),
		RWMutex: new(sync.RWMutex),
	}
}

func (p *Memory) StoreType() string {
	return InMemStore
}

func (p *Memory) IsPersistent() bool {
	return false
}

func (p *Memory) Set(name, value any) error {
	p.Lock()
	defer p.Unlock()
	p.store.Store(name, value)
	return nil
}

func (p *Memory) Get(name any) (any, bool) {
	p.RLock()
	defer p.RUnlock()
	return p.store.Load(name)
}

func (p *Memory) GetInt(name any) (int, bool) {
	return GetValue[int](p, name)
}

func (p *Memory) GetString(name any) (string, bool) {
	return GetValue[string](p, name)
}

func (p *Memory) GetBool(name any) (bool, bool) {
	return GetValue[bool](p, name)
}

func (p *Memory) GetFloat64(name any) (float64, bool) {
	return GetValue[float64](p, name)
}

func (p *Memory) Dump() map[string]interface{} {
	p.RLock()
	defer p.RUnlock()
	dumper := make(map[string]interface{})
	p.store.Range(func(key, value any) bool {
		dumper[fmt.Sprintf("%v", key)] = value
		return true
	})
	return dumper
}

func GetValue[T any](p Store, name any) (T, bool) {
	if v, ok := p.Get(name); ok {
		if rv, ok := v.(T); ok {
			return rv, true
		}
	}
	var def T
	return def, false
}

func GetObj[T any](p Store, name any) (T, error) {
	var def T
	if jsonStr, ok := GetValue[string](p, name); ok {
		return def, json.UnmarshalFromString(jsonStr, def)
	}
	return def, errors.New("not found")
}
