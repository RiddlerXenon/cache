package cache

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type Item struct {
	Value      string
	TTL        int64
	Next, Prev *Item
}

type Cache struct {
	Items   map[string]Item
	LRUList *ItemList
	Cleaner *cleaner
	Config  *config
	Mutex   sync.Mutex
}

func New() (*Cache, error) {
	cfg, err := readCacheConfig()
	if err != nil {
		return nil, err
	}

	cln := &cleaner{
		closed: make(chan bool),
	}

	c := &Cache{
		Items:   make(map[string]Item),
		LRUList: NewItemList(),
		Config:  cfg,
		Cleaner: cln,
	}

	runtime.SetFinalizer(c, closeCleaner)

	return c, nil
}

func (c *Cache) Set(key string, val string, ttl int64) {
	if ttl == 0 {
		now := time.Now()
		ttl = now.Add(time.Second * c.config.cacheConfig.TTL)
	}

	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if item, exists := c.Items[key]; exists {
		c.Items[key].Value = val
		c.Items[key].TTL = ttl
		c.LRUList.MoveToFront(item)
		return
	} else {
		item := c.LRUList.PushFront(
			&Item{
				Value: val,
				TTL:   ttl,
			},
		).Value
		c.Items[key] = item
	}

	if len(c.LRUList) == c.config.cacheConfig.Maxentries || sizeOfMap(c.Items) >= c.config.cacheConfig.MemoryLimit {
		oldest := c.LRUList.Back()
		if oldest != nil {
			c.LRUList.Remove(oldest)
			delete(c.Items, oldest)
		}
	}
}

func (c *Cache) Get(key string) (interface{}, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	item, exists := c.Items[key]
	if !exists {
		return nil, fmt.Errorf("key %s not found", key)
	}

	if item.Value == nil {
		return nil, fmt.Errorf("value for key %s is empty", key)
	}

	c.LRUList.MoveToFront(item)
	return item.Value, nil
}

func (c *Cache) DeleteExpired() {
	now := time.Now().UnixNano()

	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	for key, item := range c.Items {
		if now >= item.TTL {
			c.LRUList.Remove(item)
			delete(c.Items, key)
		}
	}
}

// func (c *Cache) CreateJSONBackup() {
// 	jsonData, err := json.MarshalIndent(c.Items, "", "  ")
// 	if err != nil {
// 		zap.S().Error(err)
// 		return
// 	}

// 	file, err := os.Create("backups/backup.json")
// 	if err != nil {
// 		zap.S().Error(err)
// 		return
// 	}
// 	defer file.Close()

// 	_, err = file.Write(jsonData)
// 	if err != nil {
// 		zap.S().Error(err)
// 		return
// 	}
// }
