package types

import (
	"context"
	"sync"
	"time"
)

type cacheValue struct {
	v          interface{}
	accessTime time.Time
	lifetime   time.Duration
}

type Cache struct {
	datas          map[uint64]*cacheValue
	mtx            sync.RWMutex
	gcInterval     time.Duration // default 1m
	lifetime       time.Duration //default 30m
	updateOnAccess bool
}

func NewCache(ctx context.Context, gc, lifeTime time.Duration, uoa bool) *Cache {
	c := &Cache{
		datas:          map[uint64]*cacheValue{},
		gcInterval:     gc,
		lifetime:       lifeTime,
		updateOnAccess: uoa,
	}
	c.Run(ctx)
	return c
}

func (c *Cache) SetDefault(k uint64, value interface{}) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	v := &cacheValue{
		v:          value,
		accessTime: time.Now(),
		lifetime:   c.lifetime,
	}

	c.datas[k] = v
}

func (c *Cache) Set(k uint64, value interface{}, lf time.Duration) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	v := &cacheValue{
		v:          value,
		accessTime: time.Now(),
		lifetime:   lf,
	}

	c.datas[k] = v
}

func (c *Cache) Get(key uint64) (interface{}, bool) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	v, ok := c.datas[key]
	if ok {
		return v.v, ok
	} else {
		return nil, ok
	}
	//false means no register yet

}

func (c *Cache) Run(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(time.Duration(c.gcInterval))
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				c.gc()
			}
		}
	}()
}

func (c *Cache) gc() {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	now := time.Now()
	for k, value := range c.datas {
		if value.accessTime.Add(time.Duration(value.lifetime)).Before(now) {
			delete(c.datas, k)
		}
	}
}
