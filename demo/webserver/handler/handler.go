package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/kawasoki/gowork/initdatabase"
	"golang.org/x/sync/singleflight"
	"log"
	"sync/atomic"
	"time"
	"ztool/database"
)

type Game struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

var c1, c2, c3 int64

func Handler(c *gin.Context) {
	//key := c.DefaultQuery("key", "default")
	//log.Println("func handler", key)
	name := GetFromRedis()
	if name == "" {
		name = GetFromDb()
		go CacheToRedis(name)
	}
	//数据库查询
	c.JSON(200, gin.H{"name": name})
}

func GetFromDb() (name string) {
	var game Game
	initdatabase.GetDB().Where("id=1").First(&game)
	//time.Sleep(5000 * time.Millisecond)
	name = game.Name
	atomic.AddInt64(&c1, 1)
	log.Println("数据库查询次数:", c1)
	return
}

func GetFromRedis() (name string) {
	name = database.GetDefaultClient().Get(context.Background(), "gameName").Val()
	//time.Sleep(50 * time.Millisecond)
	atomic.AddInt64(&c2, 1)
	//log.Println("读取缓存次数:", c2)
	return name
}

func CacheToRedis(name string) {
	atomic.AddInt64(&c3, 1)
	database.GetDefaultClient().Set(context.Background(), "gameName", name, 30*time.Second)
	log.Println("缓存数据次数:", c3)
}

var sfGroup = singleflight.Group{}
var key = "kkkk"

func Handler2(c *gin.Context) {
	name, err := loadChan(context.Background(), key)
	if err != nil {
		log.Println(err.Error())
		c.JSON(200, gin.H{"err": err.Error()})
		return
	}
	c.JSON(200, gin.H{"name": name})
}

var count int64

func loadChan(ctx context.Context, key string) (string, error) {
	data := GetFromRedis()
	if data == "" {
		atomic.AddInt64(&count, 1)
		//log.Println(count)
		// 使用 DoChan 结合 select 做超时控制
		result := sfGroup.DoChan(key, func() (interface{}, error) {
			log.Println("新来的")
			go func() {
				time.Sleep(100 * time.Millisecond)
				log.Printf("Deleting key: %v\n", key)
				sfGroup.Forget(key)
			}()

			//forgetTimer := time.AfterFunc(1000*time.Millisecond, func() {
			//	log.Println("删除kkk")
			//	sfGroup.Forget(key)
			//})
			//defer func() {
			//	log.Println("关闭定时器")
			//	forgetTimer.Stop()
			//}()
			data = GetFromDb()
			CacheToRedis(data)
			return data, nil
		})
		select {
		case r := <-result:
			return r.Val.(string), r.Err
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
	return data, nil
}
