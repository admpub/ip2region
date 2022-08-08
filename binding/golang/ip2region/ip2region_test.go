package ip2region

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

func BenchmarkMemorySearch(B *testing.B) {
	region, err := New("../../../data/ip2region.xdb")
	if err != nil {
		B.Error(err)
	}
	for i := 0; i < B.N; i++ {
		region.MemorySearch("127.0.0.1")
	}
	region.Close()
}

func TestRace(t *testing.T) {
	region, err := New("../../../data/ip2region.xdb")
	if err != nil {
		panic(err)
	}
	wg := sync.WaitGroup{}
	n := 10
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer func() {
				wg.Done()
				if e := recover(); e != nil {
					fmt.Println(`[`+strconv.Itoa(i)+`]co:`, e)
				}
			}()
			info, err := region.MemorySearch(`127.0.0.1`)
			if err != nil {
				t.Error(err)
			}
			fmt.Printf("MemorySearch: %#v\n", info)

			time.Sleep(100 * time.Millisecond)
		}(i)
	}
	for i := 0; i < 3; i++ {
		fmt.Println(`try reload:`, i)
		if err := region.Reload(); err != nil {
			t.Error(err)
		}
		time.Sleep(500 * time.Millisecond)
	}
	fmt.Println(`END`)
	wg.Wait()
	region.Close()
}
