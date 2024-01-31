package cache

import "time"

// default 超过24小时  触发 callback ，callback自定义，数据被存入数据库，再次访问时被拿出 ， status 为 -1
// expire  超过过期时间 触发 callback ，callback自定义，数据存入数据库   ，不被拿出，status 为 1
// delete  超过过期时间 触发 callback ，callback自定义，数据存入数据库   ，不被拿出，status 为 2
const (
	ITEM_STATUS_DEFAULT    = 1
	ITEM_STATUS_EXPIRE     = 2
	ITEM_STATUS_DELETE     = 4
	ITEM_STATUS_FROM_CACHE = 3
	DEFAULT_EXPIRE_TIME    = 24 * time.Hour
)

type IItem interface {
	Expired() bool
	CanExpire() bool
	SetExpireAt(t time.Time)
}

type Item struct {
	v      interface{}
	expire time.Time
	status int
}

func (i *Item) Expired() bool {
	if !i.CanExpire() {
		return false
	}
	return time.Now().After(i.expire)
}

func (i *Item) CanExpire() bool {
	return !i.expire.IsZero()
}

func (i *Item) SetExpireAt(t time.Time) {
	i.expire = t
}
