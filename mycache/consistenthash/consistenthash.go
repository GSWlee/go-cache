package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func([]byte) uint32

type Map struct {
	//哈希函数
	hash Hash
	//虚拟节点数
	relicas int
	//虚拟节点在哈希环上的位置
	key []int
	// 虚拟节点到真实节点的映射
	hashmap map[int]string
}

func New(replicas int ,hash Hash) *Map {
	m:=&Map{
		relicas: replicas,
		key: make([]int,0),
		hashmap: make(map[int]string),
		hash: crc32.ChecksumIEEE,
	}
	if hash!=nil{
		m.hash=hash
	}
	return m
}

func (m *Map) Add(keys... string)  {
	for _,key:=range keys{
		for i:=0;i<m.relicas;i++{
			hash:=int(m.hash([]byte(strconv.Itoa(i)+key)))
			m.key=append(m.key,hash)
			m.hashmap[hash]=key
		}
	}
	sort.Ints(m.key)
}

func (m *Map) Get(key string) string {
	if len(m.key)==0{
		return ""
	}

	hash:=int(m.hash([]byte(key)))
	
	id:=sort.Search(len(m.key), func(i int) bool {
		return m.key[i]>=hash
	})
	return m.hashmap[m.key[id%len(m.key)]]
}