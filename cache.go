package rbsa

import (
	"os"
	"time"
)

type (
	Cache struct {
		lock chan bool
		lookup map[string]CacheEntry
		size int
	}
	CacheEntry struct {
		access int64
		data interface{}
	}
)

func NewCache(size int) *Cache {
	c := &Cache{make(chan bool, 1),make(map[string]CacheEntry), size}
	
	c.lock<- true
	
	return c
}
func (this *Cache) Get(key string, fill func()(interface{},os.Error)) (interface{},os.Error) {
	<-this.lock
	entry, ok := this.lookup[key]
	if ok {	
		entry.access = time.Nanoseconds()
		this.lock<- true
		return entry.data, nil
	}
	this.lock<- true
	
	
	v, err := fill()
	if err != nil {
		return nil, err
	}
	<-this.lock
	this.lookup[key] = CacheEntry{time.Nanoseconds(),v}
	if len(this.lookup) > this.size {
		least := ""
		last := int64(0)
		for k, e := range this.lookup {
			if last == 0 || e.access < last {
				last = e.access
				least = k
			}
		}
		
		if least != "" {
			this.lookup[least] = CacheEntry{0,nil}, false
		}
	}
	this.lock<- true
	
	return v, nil
}
