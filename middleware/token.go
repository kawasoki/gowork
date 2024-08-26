package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
)

type MyResponse struct {
	Code int `json:"code"`
}

// AdminLog 记录日志
func AdminLog() func(c *gin.Context) {
	return func(c *gin.Context) {
		blw := &middleware.CustomResponseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		req, _ := c.GetRawData()
		reqBody := string(req)
		c.Request.Body = ioutil.NopCloser(bytes.NewBufferString(reqBody)) // Write body back

		//fmt.Println("canshu ", reqBody)

		c.Next()
		//fmt.Printf("url=%s, status=%d, resp=%s", c.Request.URL.Path, c.Writer.Status(), blw.body.String())
		//fmt.Println()
		if c.Writer.Status() != 200 {
			return
		}

		//log.Println("请求参数1", string(data))
		path := c.Request.URL.Path
		if path == "/api/v2/system/token" {
			return
		}
		var res MyResponse
		err := json.Unmarshal(blw.body.Bytes(), &res)
		if err != nil {
			fmt.Printf("url=%s, status=%d, resp=%s", c.Request.URL.Path, c.Writer.Status(), blw.body.String())
			fmt.Println()
			log.Println("记录日志参数json解析失败:", err.Error())
			return
		}
		//POST请求 并且成功就记录日志
		if c.Request.Method == "POST" && res.Code == 200 {
			adminId := c.GetInt("admin")
			go writeLog(path, adminId)
		}
	}
}
