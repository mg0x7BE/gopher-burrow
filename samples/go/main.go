//go:build !js

// Command inventory shows every Go element the theme colors.
package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
)

//go:generate stringer -type=Status

// Status is the state of an item.
type Status int

const (
	Draft Status = iota
	Published
	Archived
)

const maxItems = 1_000

var ErrNotFound = errors.New("item not found")

var defaultTags = []string{"go", "theme"}

// Item is a single inventory entry.
type Item struct {
	ID     int64          `json:"id"`
	Name   string         `json:"name,omitempty"`
	Price  float64        `json:"price"`
	Tags   []string       `json:"tags"`
	Meta   map[string]any `json:"-"`
	status Status
	notify func(string)
}

// Store keeps items safe for concurrent use.
type Store interface {
	Get(id int64) (*Item, error)
	Put(item Item) error
}

type memStore struct {
	mu    sync.RWMutex
	items map[int64]Item
}

func newStore() *memStore {
	return &memStore{items: make(map[int64]Item, maxItems)}
}

// Get returns the item or ErrNotFound.
func (s *memStore) Get(id int64) (*Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.items[id]
	if !ok {
		return nil, fmt.Errorf("get %d: %w", id, ErrNotFound)
	}
	return &item, nil
}

func (s *memStore) Put(item Item) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.items) >= maxItems {
		return errors.New("store is full")
	}
	s.items[item.ID] = item
	return nil
}

// Filter keeps the elements matching keep.
func Filter[T any](in []T, keep func(T) bool) []T {
	var out []T
	for _, v := range in {
		if keep(v) {
			out = append(out, v)
		}
	}
	return out
}

func main() {
	var store Store = newStore()
	item := Item{ID: 42, Name: "Gopher plush", Price: 19.99, Tags: defaultTags, status: Published}
	item.notify = func(msg string) { fmt.Println("notify:", msg) }

	if err := store.Put(item); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	got, err := store.Get(7)
	if errors.Is(err, ErrNotFound) {
		fmt.Printf("missing item: %v\n", err)
	}

	results := make(chan string, 2)
	var wg sync.WaitGroup
	for _, name := range []string{"alpha", "beta"} {
		wg.Go(func() {
			results <- strings.ToUpper(name)
		})
	}
	wg.Wait()
	close(results)

outer:
	for {
		select {
		case r, ok := <-results:
			if !ok {
				break outer
			}
			fmt.Println(r)
		}
	}

	for i := range 3 {
		got, err := store.Get(int64(i)) // shadows got and err from above
		if err != nil {
			continue
		}
		got.notify("found " + got.Name)
	}

	raw := `C:\path\with\backslashes`
	ratio, mask, letter := 0.75, 0xFF, 'g'
	fmt.Println(raw, ratio, mask, letter, true, got == nil)
	cheap := Filter([]Item{item}, func(i Item) bool { return i.Price < 20 })
	fmt.Println(len(cheap), item.status == Published, Draft, Archived)
}
