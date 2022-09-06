package store

type Store interface {
	StoreType() string
	IsPersistent() bool
	Set(k any, v any) error
	Get(k any) (any, bool)
	Dump() map[string]interface{}
}

var Factory = map[string]func() Store{
	InMemStore: NewMemStore,
}
