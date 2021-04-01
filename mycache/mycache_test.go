package mycache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func TestGetterFunc_Get(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		f       GetterFunc
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			f: GetterFunc(func(key string) ([]byte, error) {
				return []byte(key), nil
			}),
			args: args{key: "key"},
			want: []byte("key"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetGroup(t *testing.T) {
	var db = map[string]string{
		"Tom":  "630",
		"Jack": "589",
		"Sam":  "567",
	}

	loadCounts := make(map[string]int, len(db))
	gee := NewGroup("scores", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				if _, ok := loadCounts[key]; !ok {
					loadCounts[key] = 0
				}
				loadCounts[key] += 1
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	for k,v := range db{
		if value,err:=gee.Get(k);err!=nil||value.String()!=v{
			t.Fatalf("failed get value of %s", k)
		}

		if 	_,err:=gee.Get(k);err!=nil||loadCounts[k]>1{
			t.Fatalf("cache miss %s", k)
		}
	}

	if  value,err:=gee.Get("unKnow");err==nil{
		t.Fatalf("the value of unknow should be empty, but %s got",value)
	}
}