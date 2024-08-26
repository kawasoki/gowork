package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kawasoki/gowork/logger"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	UserId  = "user_id"
	Os      = "os"
	Channel = "channel"
	Version = "version"
	Body    = "body"
	OAid    = "oaid"
	Mid     = "mid"
	MaxSize = 1024 * 512
)

var (
	open = false
)

func SetBodyPrintOpen(o bool) {
	open = o
}
func CustomLogFormat(param gin.LogFormatterParams) string {
	//var statusColor, methodColor, resetColor string
	//if param.IsOutputColor() {
	//	statusColor = param.StatusCodeColor()
	//	methodColor = param.MethodColor()
	//	resetColor = param.ResetColor()
	//}

	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}
	BodyStr := ""
	if open {
		if i := param.Keys[Body]; i != nil {
			if raw, ok := i.([]byte); ok {
				if len(raw) != 0 {
					buffer := new(bytes.Buffer)
					if err := json.Compact(buffer, raw); err != nil {
						fmt.Println(err)
					}
					BodyStr = buffer.String()
				}
			}
		}
	}
	return fmt.Sprintf("[GIN] [%d] [%v] [%s] [%s] [%s] [version:%s] [oaid:%s] [mid:%s] [user_id:%s] [os:%s] [channel:%s] [body:%s]", param.StatusCode, param.Latency, param.ClientIP, param.Method, param.Path, param.Keys[Version], param.Keys[OAid], param.Keys[Mid], param.Keys[UserId], param.Keys[Os], param.Keys[Channel], BodyStr)
}
func CustomGinLogger(conf gin.LoggerConfig) gin.HandlerFunc {
	formatter := conf.Formatter
	if formatter == nil {
		formatter = CustomLogFormat
	}

	out := conf.Output
	if out == nil {
		out = gin.DefaultWriter
	}

	notlogged := conf.SkipPaths

	var skip map[string]struct{}

	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		c.Next()
		if _, ok := skip[path]; !ok {
			param := gin.LogFormatterParams{
				Request: c.Request,
				Keys:    c.Keys,
			}
			param.TimeStamp = time.Now()
			param.Latency = param.TimeStamp.Sub(start)

			param.ClientIP = c.ClientIP()
			param.Method = c.Request.Method
			param.StatusCode = c.Writer.Status()
			param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()

			param.BodySize = c.Writer.Size()
			param.Keys = map[string]interface{}{}
			param.Keys[UserId] = c.GetHeader(UserId)
			param.Keys[Os] = c.GetHeader(Os)
			param.Keys[Channel] = c.GetHeader(Channel)
			param.Keys[Version] = c.GetHeader(Version)
			param.Keys[OAid] = c.GetHeader(OAid)
			param.Keys[Mid] = c.GetHeader(Mid)
			if c.Keys != nil {
				param.Keys[Body] = c.Keys[Body]
			}
			if raw != "" {
				path = path + "?" + raw
			}
			param.Path = path
			logger.Info(formatter(param))
		}
	}
}

func GetHeaderOs(c *gin.Context) string {
	return c.GetHeader(Os)
}
func GetHeaderUserId(c *gin.Context) string {
	return c.GetHeader(UserId)
}

func SetBody() gin.HandlerFunc {
	return func(context *gin.Context) {
		if !open {
			return
		}
		if context.Request.Method == http.MethodPost {
			contentType := context.Request.Header.Get("Content-Type")
			if !strings.Contains(contentType, "application/json") {
				return
			}
			bodyRaw, err := io.ReadAll(context.Request.Body)
			if err != nil {
				logger.Error("SetBody err:%s", err.Error())
				return
			}
			if context.Keys == nil {
				context.Keys = map[string]interface{}{}
			}
			context.Keys["body"] = bodyRaw
			context.Request.Body = io.NopCloser(bytes.NewBuffer(bodyRaw))
		}
	}
}

func AddHeaderParam() gin.HandlerFunc {
	return func(c *gin.Context) {
		osHeader := c.GetHeader("os")
		userId := c.GetHeader("user_id")
		if c.Request.Method == http.MethodPost {
			var body map[string]interface{}
			if err := c.ShouldBindJSON(&body); err != nil {
				logger.Error("AddHeaderParam error:", err.Error())
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
				c.Abort()
				return
			}
			body["header_os"] = osHeader
			body["header_user_id"] = userId
			newBody, _ := json.Marshal(body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(newBody))
		} else {
			query := c.Request.URL.Query() // 获取现有的URL参数
			query.Add("header_os", osHeader)
			query.Add("header_user_id", userId)
			c.Request.URL.RawQuery = query.Encode() // 重新设置URL的RawQuery
		}
		c.Next()
	}
}
