package mycache

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"./consistenthash"
)

const (
	defaultBasePath = "/_cache/"
	defaultReplicas = 50
)

//核心组件，存储每个节点的基本信息
type HTTPPool struct {
	self string
	basePath string
	mu sync.Mutex
	// peers为哈希环，通过用户请求的key值来判断是由那个节点来负责处理请求
	peers *consistenthash.Map
	// httpGetter为具体实现节点缓存请求的接口，每个缓存节点对应于一个httpGetter
	httpGetters map[string]*httpGetter
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self: self,
		basePath: defaultBasePath,
	}
}

func (p *HTTPPool) Log(format string,v...interface{})  {
	log.Println("[Server %s] %s",p.self,fmt.Sprintf(format,v...))
}

func (p *HTTPPool) ServeHTTP(w  http.ResponseWriter,r *http.Request)  {

	if !strings.HasPrefix(r.URL.Path,p.basePath){
		panic("HTTPPool serving unexpected path: "+r.URL.Path)
	}
	p.Log("%s %s",r.Method,r.URL.Path)
	//最后一个参数代表最多分多少段，最后一段包含"/"

	parts:=strings.SplitN(r.URL.Path[len(p.basePath):],"/",2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	groupsname:=parts[0]
	key:=parts[1]

	g:=GetGroup(groupsname)

	if g==nil{
		http.Error(w,"no such group: "+groupsname,http.StatusNotFound)
		return
	}

	value,err:=g.Get(key)
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(value.ByteSlice())

}

type httpGetter struct {
	baseURL string
}

func (h *httpGetter) Get(group string,key string) ([]byte,error) {
	u:=fmt.Sprintf(
		"%v%v%v",
		h.baseURL,
		url.QueryEscape(group),
		url.QueryEscape(key),
		)
	res,err:=http.Get(u)
	if err!=nil{
		return nil, err
	}
	defer  res.Body.Close()
	if res.StatusCode!=http.StatusOK{
		return nil, fmt.Errorf("server return %v",res.Status)
	}

	bytes,err:=ioutil.ReadAll(res.Body)
	if err!=nil{
		return nil, fmt.Errorf("reading response body: %v",err)
	}

	return bytes, nil
}

func (p *HTTPPool) Set(peers... string)  {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers=consistenthash.New(defaultReplicas,nil)
	p.peers.Add(peers...)
	p.httpGetters=make(map[string]*httpGetter)
	for _,peer:=range peers{
		p.httpGetters[peer]=&httpGetter{baseURL: peer+p.basePath}
	}
}

func (p *HTTPPool) PickPeer(key string) (PeerGetter,bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if peer:=p.peers.Get(key);peer!=""&&peer!=p.self{
		p.Log("Pick peer %s",peer)
		return p.httpGetters[peer],true
	}
	return nil,false
}
