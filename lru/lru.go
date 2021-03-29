package lru

import "container/list"

type Cache struct {
	maxBytes  int64
	nBytes    int64
	ll        *list.List
	cache     map[string]*list.Element
	OnEvicted func(key string, value Value)
}

type entry struct {
	key string
	value Value
}

// key-Value对中的Value,为简化只用Len()来表示占用多少字节
type Value interface {
	Len() int
}

func New(maxBytes int64,OnEvicted func(key string,value Value)) *Cache {
	return &Cache{
		maxBytes: maxBytes,
		ll: list.New(),
		cache: make(map[string]*list.Element),
		OnEvicted: OnEvicted,
	}
}

//访问cache
func (c *Cache)Get(key string) (value Value,ok bool) {
	if ele,ok:=c.cache[key];ok{
		c.ll.MoveToFront(ele)
		//x.(type) 类型断言，x通常为interface,判断x是否为type如果是则转换成type
		kv:=ele.Value.(*entry)
		return kv.value,true
	}
	return
}

//删除
func (c *Cache)DeleteOld()  {
	ele:=c.ll.Back()
	if ele!=nil{
		c.ll.Remove(ele)
		kv:=ele.Value.(*entry)
		delete(c.cache,kv.key)
		c.nBytes=c.nBytes-int64(len(kv.key))-int64(kv.value.Len())
		if c.OnEvicted!=nil{
			c.OnEvicted(kv.key,kv.value)
		}
	}
}

//添加/修改元素
func (c *Cache) Add(key string,value Value)  {
	if ele,ok:=c.cache[key];ok{
		c.ll.MoveToFront(ele)
		kv:=ele.Value.(*entry)
		ele.Value=&entry{key: key,value: value}
		c.nBytes=c.nBytes-int64(kv.value.Len())+int64(value.Len())
	}else{
		kv:=&entry{key: key,value: value}
		c.ll.PushFront(kv)
		c.cache[key]=c.ll.Front()
		c.nBytes+=int64(len(kv.key))+int64(kv.value.Len())
	}
	for c.maxBytes!=0&&c.maxBytes<c.nBytes{
		c.DeleteOld()
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
