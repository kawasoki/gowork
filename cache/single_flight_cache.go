package cache

import (
	"context"
	"errors"
	"golang.org/x/sync/singleflight"
	"time"
)

var ErrCacheMiss = errors.New("cache miss")
var g singleflight.Group

type LoadDataFunc func() (interface{}, error)

// redis层没有做siglefligt
func GetData(ctx context.Context, key string, LoadFromCache, LoadFormDb LoadDataFunc) (interface{}, error) {
	data, err := LoadFromCache()
	if err != nil && errors.Is(err, ErrCacheMiss) {
		// 使用 DoChan 结合 select 做超时控制
		result := g.DoChan(key, func() (interface{}, error) {
			go func() {
				time.Sleep(100 * time.Millisecond)
				g.Forget(key)
			}()
			return LoadFormDb()
		})
		select {
		case r := <-result:
			return r.Val, r.Err
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	return data, err
}
