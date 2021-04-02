package mycache

import (
	"fmt"
	"log"
	"sync"
	"./singleflight"
)

type Getter interface {
	Get(key string) ([]byte,error)
}

type GetterFunc func(key string)([]byte,error)

func (f GetterFunc) Get(key string) ([]byte,error) {
	return f(key)
}

type Group struct {
	name string
	getter Getter
	cache Cache
	peers PeerPicker
	loader *singleflight.Group
}

var (
	mu sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string,cacheBytes int64,getter Getter) *Group {
	if getter==nil{
		panic("nil Getter")
	}
	g:=&Group{
		name: name,
		getter: getter,
		cache: Cache{maxBytes: cacheBytes},
		loader: &singleflight.Group{},
	}
	groups[name]=g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	g:=groups[name]
	return g
}

func (g *Group) Get(key string) (Byteviews,error) {
	if key==""{
		return Byteviews{},fmt.Errorf("key is nil")
	}

	if value,ok:=g.cache.Get(key);ok{
		log.Println("[cache] hit!")
		return value,nil
	}
	return g.load(key)
}

func (g *Group) load(key string) (value Byteviews,err error) {
	viewi,err:=g.loader.Do(key, func() (interface{}, error) {
		if g.peers!=nil{
			if peer,ok:=g.peers.PickPeer(key);ok{
				if value,err=g.getFromPeer(peer,key);err==nil{
					return value,nil
				}
				log.Println("[cache] Failed to get from peer",err)
			}
		}
		return g.getLocally(key)
	})
	if err==nil{
		return viewi.(Byteviews),nil
	}
	return
}

func (g *Group) getLocally(key string) (Byteviews,error) {
	bytes,err:=g.getter.Get(key)
	if err!=nil{
		return Byteviews{}, err
	}
	value:=Byteviews{b: bytes}
	g.populateCache(key,value)
	return value,nil
}

func (g *Group) getFromPeer(peer PeerGetter,key string) (Byteviews,error){
	bytes,err:=peer.Get(g.name,key)
	if err!=nil{
		return Byteviews{},err
	}
	return Byteviews{b: bytes}, err
}

func (g *Group) populateCache(key string,value Byteviews)  {
	g.cache.Add(key,value)
}

func (g *Group) RegisterPeers(peers PeerPicker)  {
	if g.peers!=nil{
		panic("RegisterPeerPicker called more than once")
	}
	g.peers=peers
}

