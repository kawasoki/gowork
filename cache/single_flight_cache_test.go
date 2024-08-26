package cache

import (
	"sync"
	"testing"
	"time"
)

func BenchmarkGetData(b *testing.B) {

}

func TestName(t *testing.T) {
	//ctx := context.Background()
	var wg sync.WaitGroup
	count := 13
	wg.Add(count)
	for i := 0; i < count; i++ {
		if i%7 == 0 {
			time.Sleep(200 * time.Millisecond)
		}
		go func() {
			defer wg.Done()
			//GetDataV2(ctx, "str", loadFromCache, loadFromDb)
		}()
	}
	wg.Wait()
}

func loadFromCache() (interface{}, error) {

	return nil, ErrCacheMiss
}
func loadFromDb() (interface{}, error) {
	time.Sleep(2 * time.Second)
	return nil, nil
}
