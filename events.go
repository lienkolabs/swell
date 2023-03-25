package swell

import (
	"sync"

	"github.com/lienkolabs/swell/crypto"
	"github.com/lienkolabs/swell/util"
)

// Any valid transaction must start with Version on byte[0] and have the clock as
// measured in blocks since genesis set on the following 8 bytes. Encoding is
// little-endian.
type Event []byte

func (t Event) Clock() uint64 {
	if t[0] != Version || len(t) < 9 {
		return 0
	}
	clock, _ := util.ParseUint64(t, 1)
	return clock
}

func (t Event) Hash() crypto.Hash {
	return crypto.Hasher(t)
}

type Events []Event

type EventsPool struct {
	queue  []crypto.Hash // order in which instructions are received
	events map[crypto.Hash]Event
	mu     sync.Mutex
}

func NewInstructionPool() *EventsPool {
	return &EventsPool{
		queue:  make([]crypto.Hash, 0),
		events: make(map[crypto.Hash]Event),
	}
}

func (pool *EventsPool) Unqueue() (Event, crypto.Hash) {
	if len(pool.queue) == 0 {
		return nil, crypto.ZeroHash
	}
	pool.mu.Lock()
	defer pool.mu.Unlock()
	for n, hash := range pool.queue {
		if event, ok := pool.events[hash]; ok {
			pool.queue = pool.queue[n+1:]
			delete(pool.events, hash)
			return event, hash
		}
	}
	pool.queue = pool.queue[:0]
	return nil, crypto.ZeroHash
}

func (pool *EventsPool) Queue(event Event, hash crypto.Hash) {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	pool.queue = append(pool.queue, hash)
	pool.events[hash] = event
}

func (pool *EventsPool) Delete(hash crypto.Hash) {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	delete(pool.events, hash)
}

func (pool *EventsPool) DeleteArray(hashes []crypto.Hash) {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	for _, hash := range hashes {
		delete(pool.events, hash)
	}
}
