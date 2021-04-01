package main

import (
	"./mycache"
	"fmt"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	mycache.NewGroup("scores", 2<<10, mycache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
	addr:="localhost:8000"
	pools:=mycache.NewHTTPPool(addr)
	log.Println("cache is running at",addr)
	log.Fatal(http.ListenAndServe(addr,pools))
}
