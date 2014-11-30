package syncmap

import (
	"encoding/json"
	"sync"
)

const (
	defaultShardCount uint8 = 32
)

type syncMap struct {
	items map[string]interface{}
	sync.RWMutex
}

type SyncMap struct {
	shardCount uint8
	shards     []*syncMap
}

func New() *SyncMap {
	return NewWithShard(defaultShardCount)
}

func NewWithShard(shardCount uint8) *SyncMap {
	m := new(SyncMap)
	m.shardCount = shardCount
	m.shards = make([]*syncMap, m.shardCount)
	for i, _ := range m.shards {
		m.shards[i] = &syncMap{items: make(map[string]interface{})}
	}
	return m
}

func (m *SyncMap) locate(key string) *syncMap {
	return m.shards[BkdrHash(key)&uint32((m.shardCount-1))]
}

func (m *SyncMap) Get(key string) (value interface{}, ok bool) {
	shard := m.locate(key)
	shard.RLock()
	defer shard.RUnlock()

	value, ok = shard.items[key]
	return
}

func (m *SyncMap) Set(key string, value interface{}) {
	shard := m.locate(key)
	shard.Lock()
	defer shard.Unlock()
	shard.items[key] = value
}

func (m *SyncMap) Delete(key string) {
	shard := m.locate(key)
	shard.Lock()
	defer shard.Unlock()
	delete(shard.items, key)
}

func (m *SyncMap) Has(key string) bool {
	_, ok := m.Get(key)
	return ok
}

func (m *SyncMap) Size() int {
	size := 0
	for _, shard := range m.shards {
		shard.RLock()
		size += len(shard.items)
		shard.RUnlock()
	}
	return size
}

func (m *SyncMap) Flush() int {
	size := 0
	for _, shard := range m.shards {
		shard.Lock()
		size += len(shard.items)
		shard.items = make(map[string]interface{})
		shard.Unlock()
	}
	return size
}

func (m *SyncMap) IterKeys() <-chan string {
	ch := make(chan string)
	go func() {
		for _, shard := range m.shards {
			shard.RLock()
			for key, _ := range shard.items {
				ch <- key
			}
			shard.RUnlock()
		}
		close(ch)
	}()
	return ch
}

type Item struct {
	Key   string
	Value interface{}
}

func (m *SyncMap) IterItems() <-chan Item {
	ch := make(chan Item)
	go func() {
		for _, shard := range m.shards {
			shard.RLock()
			for key, value := range shard.items {
				ch <- Item{key, value}
			}
			shard.RUnlock()
		}
		close(ch)
	}()
	return ch
}

func (m SyncMap) MarshalJSON() ([]byte, error) {
	x := make(map[string]interface{})
	for item := range m.IterItems() {
		x[item.Key] = item.Value
	}
	return json.Marshal(x)
}
