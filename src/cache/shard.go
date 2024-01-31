package cache

import (
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
)

// ExpiredCallback Callback the function when the key-value pair expires
// Note that it is executed after expiration
type ExpiredCallback func(k string, v interface{}) error

type SimpCacheItem struct {
	Key   string      `json:"k" db:"k"`
	Table string      `json:"t" db:"t"`
	Value interface{} `json:"v" db:"v"`
}

func (s *SimpCacheItem) GetFromTable(db *sqlx.DB) {

}

type GetWhenExpiredFunc func(k string) (value interface{}, exist bool)
type memCacheShard struct {
	hashmap         map[string]Item
	lock            sync.RWMutex
	expiredCallback ExpiredCallback // 普通过期
	deleteCallback  ExpiredCallback // 被删除
	defaultCallback ExpiredCallback // 入库
	getWhenExpire   GetWhenExpiredFunc
}

func newMemCacheShard(conf *Config) *memCacheShard {
	return &memCacheShard{
		expiredCallback: conf.expiredCallback,
		deleteCallback:  conf.deleteCallback,
		hashmap:         map[string]Item{},
	}
}

// 没有过期时间的默认给24小时自动过期
// 防止太多cache存入内存中
func (c *memCacheShard) set(k string, item *Item) {
	if !item.CanExpire() || item.status == ITEM_STATUS_DEFAULT || item.status == ITEM_STATUS_FROM_CACHE {
		item.status = ITEM_STATUS_DEFAULT
		item.SetExpireAt(time.Now().Add(DEFAULT_EXPIRE_TIME))
	} else {
		item.status = ITEM_STATUS_EXPIRE
	}
	c.lock.Lock()
	c.hashmap[k] = *item
	c.lock.Unlock()
	return
}

func (c *memCacheShard) get(k string) (interface{}, bool) {
	c.lock.RLock()
	item, exist := c.hashmap[k]
	c.lock.RUnlock()
	if !exist {
		return nil, false
	}
	if !item.Expired() {
		return item.v, true
	}
	value, exist := c.getWhenExpire(k)
	if exist && item.status == ITEM_STATUS_DEFAULT {
		c.set(k, &Item{v: value, status: ITEM_STATUS_FROM_CACHE})
		return value, true
	}
	if c.delExpired(k) {
		return nil, false
	}
	return c.get(k)
}

func (c *memCacheShard) getSet(k string, item *Item) (interface{}, bool) {
	defer c.set(k, item)
	return c.get(k)
}

func (c *memCacheShard) getDel(k string) (interface{}, bool) {
	defer c.del(k)
	return c.get(k)
}

func (c *memCacheShard) del(k string) int {
	var count int
	c.lock.Lock()
	v, found := c.hashmap[k]
	if found {
		delete(c.hashmap, k)
		if !v.Expired() {
			count++
			c.deleteCallback(k, v)
		}
	}
	c.lock.Unlock()
	return count
}

// delExpired Only delete when key expires
func (c *memCacheShard) delExpired(k string) bool {
	c.lock.Lock()
	item, found := c.hashmap[k]
	if !found || !item.Expired() {
		c.lock.Unlock()
		return false
	}
	delete(c.hashmap, k)
	c.lock.Unlock()
	if c.expiredCallback != nil {
		switch item.status {
		case ITEM_STATUS_DEFAULT:
			{
				_ = c.defaultCallback(k, item.v)
			}
		case ITEM_STATUS_EXPIRE:
			{
				_ = c.expiredCallback(k, item.v)
			}
		}
	}
	return true
}

func (c *memCacheShard) ttl(k string) (time.Duration, bool) {
	c.lock.RLock()
	v, found := c.hashmap[k]
	c.lock.RUnlock()
	if !found || !v.CanExpire() || v.Expired() {
		return 0, false
	}
	return v.expire.Sub(time.Now()), true
}

func (c *memCacheShard) checkExpire() {
	var expiredKeys []string
	c.lock.RLock()
	for k, item := range c.hashmap {
		if item.Expired() {
			expiredKeys = append(expiredKeys, k)
		}
	}
	c.lock.RUnlock()
	for _, k := range expiredKeys {
		c.delExpired(k)
	}
}
