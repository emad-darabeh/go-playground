package main

import (
	"fmt"
	"sync"
)

type ConfigStore[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

func NewConfigStore[K comparable, V any]() *ConfigStore[K, V] {
	return &ConfigStore[K, V]{
		data: make(map[K]V),
	}
}

func (s *ConfigStore[K, V]) Set(key K, value V) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

func (s *ConfigStore[K, V]) Get(key K) (V, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	return val, ok
}

type StringConfig = ConfigStore[string, any]

func GetTyped[T any](c *StringConfig, key string) (*T, bool) {
	val, ok := c.Get(key)
	if !ok {
		return nil, false
	}

	typedVal, valid := val.(*T)
	return typedVal, valid
}

func main() {
	config := NewConfigStore[string, any]()

	config.Set("key-int", new(5))
	config.Set("key-str", new("hello"))
	config.Set("key-bool", new(true))

	if val, ok := GetTyped[int](config, "key-int"); ok {
		fmt.Printf("key %s has value %d\n", "key-int", *val)
	}

	if val, ok := GetTyped[string](config, "key-str"); ok {
		fmt.Printf("key %s has value %s\n", "key-str", *val)
	}

	if val, ok := GetTyped[bool](config, "key-bool"); ok {
		fmt.Printf("key %s has value %v\n", "key-bool", *val)
	}
}
