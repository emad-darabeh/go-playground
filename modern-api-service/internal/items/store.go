package items

import "sync"

type InMemoryStore struct {
	mu    sync.RWMutex
	items map[uint64]Item
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{items: make(map[uint64]Item)}
}

func (s *InMemoryStore) Get(id uint64) (Item, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.items[id]
	return item, ok
}

func (s *InMemoryStore) Add(item Item) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[item.ID] = item
}

func (s *InMemoryStore) Delete(id uint64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.items[id]
	delete(s.items, id)
	return ok
}

func (s *InMemoryStore) List() []Item {
	s.mu.RLock()
	defer s.mu.RUnlock()
	all := make([]Item, 0, len(s.items))
	for _, item := range s.items {
		all = append(all, item)
	}
	return all
}
