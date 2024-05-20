package data

import (
	"sync"

	"github.com/DaanV2/f1-game-dashboards/server/sessions"
)

type (
	MemoryStorage struct {
		chairs *TypedStorage[sessions.Chair]
		config *memStorage
	}

	memStorage struct {
		lock  sync.Mutex
		items map[string][]byte
	}
)

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		config: newMStorage(),
		chairs: NewTypedStorage[sessions.Chair](newMStorage()),
	}
}

func (fs *MemoryStorage) Chairs() Storage[sessions.Chair] {
	return fs.chairs
}

func (fs *MemoryStorage) Config() RawStorage {
	return fs.config
}

func newMStorage() *memStorage {
	return &memStorage{
		lock:  sync.Mutex{},
		items: make(map[string][]byte),
	}
}

func (ms *memStorage) Get(id string) ([]byte, error) {
	ms.lock.Lock()
	defer ms.lock.Unlock()

	data, ok := ms.items[id]
	if !ok {
		return data, ErrNotFound
	}

	return data, nil
}

func (ms *memStorage) Set(id string, value []byte) error {
	ms.lock.Lock()
	defer ms.lock.Unlock()

	ms.items[id] = value

	return nil
}

func (ms *memStorage) Delete(id string) error {
	ms.lock.Lock()
	defer ms.lock.Unlock()

	delete(ms.items, id)

	return nil
}

func (ms *memStorage) Keys() []string {
	ms.lock.Lock()
	defer ms.lock.Unlock()

	keys := make([]string, 0, len(ms.items))
	for k := range ms.items {
		keys = append(keys, k)
	}

	return keys
}
