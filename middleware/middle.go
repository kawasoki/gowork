package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func AccessTimeCheck(timeArgs ...[2]time.Time) func(c *gin.Context) {
	return func(c *gin.Context) {
		now := time.Now()
		for _, timeA := range timeArgs {
			if now.After(timeA[0]) && now.Before(timeA[1]) {
				c.Next()
				return
			}
		}
		c.Abort()
		c.JSON(200, gin.H{"success": false, "msg": "Access Deny", "code": 404})
	}
}

// MaxAllowed 限流
func MaxAllowed(n int) func(c *gin.Context) {
	sem := make(chan struct{}, n)
	acquire := func() { sem <- struct{}{} }
	release := func() { <-sem }
	return func(c *gin.Context) {
		acquire()
		defer release()
		c.Next()
	}
}

// ReSubmit redis做的简易防重复提交
func ReSubmit() func(c *gin.Context) {
	return func(c *gin.Context) {
		uid, ok := c.Get("user")
		if !ok {
			c.Abort()
			c.JSON(200, gin.H{"code": 400, "msg": "无"})
			return
		}
		key := fmt.Sprintf("%s_%s", uid.(string), c.Request.URL.Path)
		//redis出错setSucceed默认false 阻止继续执行
		setSucceed, _ := initialize.RedisClient.SetNX(key, 1, 5*time.Second).Result()
		if !setSucceed {
			c.Abort()
			c.JSON(200, gin.H{"success": false, "msg": "重复提交", "code": 404})
			return
		}
		c.Next()
		initialize.RedisClient.Del(key)
	}
}
