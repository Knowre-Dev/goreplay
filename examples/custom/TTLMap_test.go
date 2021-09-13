package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTTLMap(t *testing.T) {
	m := NewTTLMap(1)
	for i := 0; i < 10000; i++ {
		k, v := fmt.Sprint("key", i), fmt.Sprint("value", i)
		m.Put(k, v)
	}
	assert.Equal(t, 10000, m.Len())

	time.Sleep(5 * time.Second)
	assert.Equal(t, 0, m.Len())
	fmt.Println("len(m) (this will be empty):", m.Len())
}
