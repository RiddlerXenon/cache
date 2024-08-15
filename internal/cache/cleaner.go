package cache

import "time"

type cleaner struct {
	closed chan bool
}

func closeCleaner(c *Cache) {
	c.Cleaner.closed <- true
}

func (cln *cleaner) Run(c *Cache) {
	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case <-ticker.C:
			c.DeleteExpired()
			// c.CreateJSONBackup()
		case <-cln.closed:
			ticker.Stop()
			return
		}
	}
}
