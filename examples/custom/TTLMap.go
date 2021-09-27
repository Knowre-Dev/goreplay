package main

import (
	"sync"
	"time"
)

type item struct {
	value      string
	lastAccess int64
}

type TTLMap struct {
	sync.Mutex
	m map[string]*item
}

func NewTTLMap(config AppConfig) *TTLMap {
	m := &TTLMap{m: make(map[string]*item)}
	go func() {
		for now := range time.Tick(time.Second) {
			m.Lock()
			for k, v := range m.m {
				if now.Unix()-v.lastAccess > int64(config.MaxTTL) {
					delete(m.m, k)
				}
			}
			m.Unlock()
		}
	}()
	return m
}

func (m *TTLMap) Len() int {
	return len(m.m)
}

func (m *TTLMap) Set(k string, v string) {
	m.Lock()
	it, ok := m.m[k]
	if ok {
		m.m[k] = nil
	}

	it = &item{value: v}
	m.m[k] = it

	it.lastAccess = time.Now().Unix()
	m.Unlock()
}

func (m *TTLMap) Get(k string) (v string, exist bool) {
	m.Lock()
	if it, ok := m.m[k]; ok {
		v = it.value
		exist = true
		it.lastAccess = time.Now().Unix()
	}
	m.Unlock()
	return

}
