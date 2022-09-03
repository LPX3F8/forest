package store

import "testing"

func TestMemStore(t *testing.T) {
	s := NewMemStore()
	if err := s.Set("test", 123); err != nil {
		t.Fatal(err)
	}
	if v, ok := s.Get("test"); v != 123 || !ok {
		t.Fatal("get value failed")
	}
}
