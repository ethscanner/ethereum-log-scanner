package cache

import (
	"fmt"
	"testing"
)

func TestCache(t *testing.T) {

	c1 := NewCacheByName("1", 66*2, nil)
	c1.Add("0x10b55553238d25079beb3c3b2840a55bed27fe6c7bc4bc4f9fb1f31e7a7a9b0a", ByteView{b: []byte{}})
	c1.Add("0x10b55553238d25079beb3c3b2840a55bed27fe6c7bc4bc4f9fb1f31e7a7a9b0b", ByteView{b: []byte{}})
	c1.Add("0x10b55553238d25079beb3c3b2840a55bed27fe6c7bc4bc4f9fb1f31e7a7a9b0c", ByteView{b: []byte{}})
	b, ok := c1.Get("0x10b55553238d25079beb3c3b2840a55bed27fe6c7bc4bc4f9fb1f31e7a7a9b0b")
	fmt.Println(b.ByteSlice(), ok)

}

func TestAbc(t *testing.T) {

}
