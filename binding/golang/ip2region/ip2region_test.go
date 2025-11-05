package ip2region

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/admpub/ip2region/v3/binding/golang/xdb"
)

func BenchmarkMemorySearch(B *testing.B) {
	region, err := New("../../../data/ip2region.xdb", true)
	if err != nil {
		B.Error(err)
	}
	for i := 0; i < B.N; i++ {
		region.MemorySearch("127.0.0.1")
	}
	region.Close()
}

func TestRace(t *testing.T) {
	region, err := New("../../../data/ip2region_v6.xdb", false)
	if err != nil {
		panic(err)
	}
	testIP := `219.133.111.87`
	if region.dbVer.Id == xdb.IPv6VersionNo {
		testIP = `240e:87c:892:ffff:ffff:ffff:ffff:ffff`
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
			result, err := region.MemorySearchString(testIP)
			if err != nil {
				t.Error(err)
			}
			info := ParseResult(result)
			fmt.Printf("MemorySearch: %#v [%s]\n", info, result)

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
