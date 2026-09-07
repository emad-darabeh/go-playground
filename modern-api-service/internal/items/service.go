package items

import (
	"errors"
	"math/rand/v2"
)

var (
	ErrNotFound       = errors.New("item not found")
	ErrNameRequired   = errors.New("name is required")
)

type Item struct {
	ID   uint64
	Name string
}

type Store interface {
	Get(id uint64) (Item, bool)
	Add(item Item)
	Delete(id uint64) bool
	List() []Item
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) Create(name string) (Item, error) {
	if name == "" {
		return Item{}, ErrNameRequired
	}
	item := Item{ID: rand.Uint64(), Name: name}
	s.store.Add(item)
	return item, nil
}

func (s *Service) Get(id uint64) (Item, error) {
	item, ok := s.store.Get(id)
	if !ok {
		return Item{}, ErrNotFound
	}
	return item, nil
}

func (s *Service) Delete(id uint64) error {
	if !s.store.Delete(id) {
		return ErrNotFound
	}
	return nil
}

func (s *Service) List() []Item {
	return s.store.List()
}
