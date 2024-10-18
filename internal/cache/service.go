package cache

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type Item struct {
	Key        string
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

	go cln.Run(c)
	runtime.SetFinalizer(c, closeCleaner)

	return c, nil
}

func (c *Cache) Set(key string, val string, ttl int64) {
	if ttl == 0 {
		now := time.Now()
		ttl = now.Add(time.Second * time.Duration(c.Config.Cache.TTL)).UnixNano()
	}

	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if item, exists := c.Items[key]; exists {
		item.Value = val
		item.TTL = ttl
		c.LRUList.MoveToFront(&item)
		return
	} else {
		item := c.LRUList.PushFront(
			&Item{
				Key:   key,
				Value: val,
				TTL:   ttl,
			},
		)
		c.Items[key] = *item
	}

	if c.LRUList.Len == c.Config.Cache.MaxEntries || sizeOfMap(c.Items) >= c.Config.Cache.MemoryLimit {
		c.DeleteExpired()

		oldest := c.LRUList.Back()
		if oldest != nil {
			c.LRUList.Remove(oldest)
			delete(c.Items, oldest.Key)
		}
	}
}

func (c *Cache) Get(key string) (string, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	item, exists := c.Items[key]
	if !exists {
		return "", fmt.Errorf("key %s not found", key)
	}

	if item.Value == "" {
		return "", fmt.Errorf("value for key %s is empty", key)
	}

	c.LRUList.MoveToFront(&item)
	return item.Value, nil
}

func (c *Cache) DeleteExpired() {
	now := time.Now().UnixNano()

	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	for key, item := range c.Items {
		if now >= item.TTL {
			c.LRUList.Remove(&item)
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
