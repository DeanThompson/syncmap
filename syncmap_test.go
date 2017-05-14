package syncmap

import (
	"sort"
	"strconv"
	"testing"
)

func Test_New(t *testing.T) {
	m1 := New()
	if m1 == nil {
		t.Error("New(): map is nil")
	}
	if m1.shardCount != defaultShardCount {
		t.Error("New(): map's shard count is wrong")
	}
	if m1.Size() != 0 {
		t.Error("New(): new map should be empty")
	}

	var shardCount uint8 = 64
	m2 := NewWithShard(shardCount)
	if m2 == nil {
		t.Error("NewWithShard(): map is nil")
	}
	if m2.shardCount != shardCount {
		t.Error("NewWithShard(): map's shard count is wrong")
	}
	if m2.Size() != 0 {
		t.Error("New(): new map should be empty")
	}
}

func Test_Set(t *testing.T) {
	m := New()
	m.Set("one", 1)
	m.Set("tow", 2)
	if m.Size() != 2 {
		t.Error("map should have 2 items.")
	}
}

func Test_Get(t *testing.T) {
	m := New()
	v1, ok := m.Get("not_exist_at_all")
	if ok {
		t.Error("ok should be false when key is missing")
	}
	if v1 != nil {
		t.Error("value should be nil for missing key")
	}

	m.Set("one", 1)

	v2, ok := m.Get("one")
	if !ok {
		t.Error("ok should be true when key exists")
	}
	if 1 != v2.(int) {
		t.Error("value should be an integer of value 1")
	}
}

func Test_Has(t *testing.T) {
	m := New()
	if m.Has("missing_key") {
		t.Error("Has should return False for missing key")
	}

	m.Set("one", 1)
	if !m.Has("one") {
		t.Error("Has should return True for existing key")
	}
}

func Test_Delete(t *testing.T) {
	m := New()
	m.Set("one", 1)
	m.Delete("one")
	if m.Has("one") {
		t.Error("Delete shoudl remove the given key from map")
	}
}

func Test_Size(t *testing.T) {
	m := New()
	for i := 0; i < 42; i++ {
		m.Set(strconv.Itoa(i), i)
	}
	if m.Size() != 42 {
		t.Error("Size doesn't return the right number of items")
	}
}

func Test_Flush(t *testing.T) {
	var shardCount uint8 = 64
	m := NewWithShard(shardCount)
	for i := 0; i < 42; i++ {
		m.Set(strconv.Itoa(i), i)
	}
	count := m.Flush()
	if count != 42 {
		t.Error("Flush should return the size before removing")
	}
	if m.Size() != 0 {
		t.Error("Flush should remove all items from map", m.Size())
	}
	if m.shardCount != shardCount {
		t.Error("map should have the same shardCount after Flush")
	}
}

func Test_IterKeys(t *testing.T) {
	loop := 100
	expectedKeys := make([]string, loop)

	m := New()
	for i := 0; i < loop; i++ {
		key := strconv.Itoa(i)
		expectedKeys[i] = key
		m.Set(key, i)
	}

	keys := make([]string, 0)
	for key := range m.IterKeys() {
		keys = append(keys, key)
	}

	if len(keys) != len(expectedKeys) {
		t.Error("IterKeys doesn't loop the right times")
	}

	sort.Strings(keys)
	sort.Strings(expectedKeys)

	for i, v := range keys {
		if v != expectedKeys[i] {
			t.Error("IterKeys doesn't loop over the right keys")
		}
	}
}

func Test_Pop(t *testing.T) {
	m := New()
	// m.Pop()

	m.Set("one", 1)

	k, v := m.Pop()
	if k != "one" && v.(int) != 1 {
		t.Error("Pop should returns the only item")
	}
	if m.Size() != 0 {
		t.Error("Size should be 0 after pop the only item")
	}
}

func Test_ForEachItem(t *testing.T) {
	m := New()
	for i := 0; i < 42; i++ {
		m.Set(strconv.Itoa(i), i)
	}
	var i *Item
	m.ForEachItem(func(item *Item) bool {
		i = item
		return false
	})
	m.Set(i.Key, "new")
	if v, ok := m.Get(i.Key); !ok || v.(string) != "new" {
		t.Error("value should be an string of value new")
	}

}

func Test_ForEachKey(t *testing.T) {
	m := New()
	for i := 0; i < 42; i++ {
		m.Set(strconv.Itoa(i), i)
	}
	var k string
	m.ForEachKey(func(key string) bool {
		k = key
		return false
	})
	m.Delete(k)
	v1, ok := m.Get(k)
	if ok {
		t.Error("ok should be false when key is missing")
	}
	if v1 != nil {
		t.Error("value should be nil for missing key")
	}
}
