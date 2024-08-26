package main

import (
	"encoding/json"
	"fmt"
	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
	"github.com/kawasoki/gowork/demo/webserver/handler"
	"log"
	"net/http"
	"time"
)

func main() {
	r := gin.New()
	r.Use(gin.Recovery())
	gin.SetMode(gin.ReleaseMode)
	memoryStore := persist.NewMemoryStore(1 * time.Minute)
	r.GET("/game", cache.CacheByRequestURI(
		memoryStore,
		2*time.Second,
		//cache.WithOnHitCache(mid2),
		//cache.WithOnShareSingleFlight(mid3),
	), handler.Handler)
	r.GET("/game2", handler.Handler2)
	// 设置路由
	r.GET("/events", Cors(), checkOrderSuccessHandler)
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}

func mid2(c *gin.Context) {
	fmt.Println("func mid2")
}

func mid3(c *gin.Context) {
	fmt.Println("func mid3")

}

func checkOrderSuccessHandler(c *gin.Context) {
	order, _ := c.GetQuery("order")
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	var res = SSEResult{Code: http.StatusOK, Info: "init"}
	//如果
	// 发送初始事件
	fmt.Fprintf(c.Writer, formatStr(res))

	var ch = make(chan struct{})
	go orderBack(ch, order)

	closeNotify := c.Writer.CloseNotify()
	// 定义发送的事件数量上限
	maxEvents := 10
	// 向客户端发送事件数据

	for i := 0; i < maxEvents; i++ {
		select {
		case <-ch:
			res.Code = http.StatusOK
			res.Info = "已付款"
			c.Writer.WriteString(formatStr(res))
			c.Writer.Flush()
			fmt.Println("已支付")
			return
		case <-closeNotify:
			fmt.Println("客户端断开") // 客户端连接已关闭，退出
			return
		default:
			// 格式化事件数据，每隔一秒钟发送一个事件
			res.Code = http.StatusNoContent
			res.Info = time.Now().Format(time.DateTime)
			c.Writer.WriteString(formatStr(res))
			c.Writer.Flush() // 刷新缓冲区，确保数据被立即发送到客户端
			log.Println(i)
			// 等待一秒钟
			<-time.After(1 * time.Second)
		}
	}
	res.Code = http.StatusOK
	res.Info = "超时断开"
	fmt.Println("超时断开")
	c.Writer.WriteString(formatStr(res))
	c.Writer.Flush()
}

type SSEResult struct {
	Code                int    `json:"code"`
	Info                string `json:"info"`
	ProductId           int    `json:"product_id"`
	StudentVerifyStatus int    `json:"student_verify_status"`
}

func formatStr(s SSEResult) string {
	b, _ := json.Marshal(s)
	return fmt.Sprintf("data: %s\n\n", string(b))
}

// 3秒后回调信息
func orderBack(ch chan struct{}, orderNumber string) {
	time.Sleep(4 * time.Second)
	ch <- struct{}{}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token,Authorization,Token,x-requested-with")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}
