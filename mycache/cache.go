package mycache

import (
	"./lru"
	"sync"
)

type Cache struct {
	mu sync.Mutex
	c *lru.Cache
	maxBytes int64
}

func (c *Cache) Add(key string,value Byteviews)  {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.c==nil{
		c.c=lru.New(c.maxBytes,nil)
	}
	c.c.Add(key,value)
}

func (c *Cache) Get(key string)  (value Byteviews,ok bool){
	if c.c==nil{
		return
	}else{
		c.mu.Lock()
		defer c.mu.Unlock()
		if value,ok:=c.c.Get(key);ok{
			return value.(Byteviews),ok
		}
	}
	return
}