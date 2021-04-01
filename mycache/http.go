package mycache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/_cache/"

type HTTPPool struct {
	self string
	basePath string
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