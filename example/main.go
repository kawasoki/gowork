package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kawasoki/gowork/gconf"
	"github.com/kawasoki/gowork/middleware"
	"github.com/kawasoki/gowork/response"
	"github.com/kawasoki/gowork/server_model"
	"log"
)

var validate *validator.Validate

// 初始化Gin引擎
func main() {
	//gconf.LoadGConf("tm")
	gin.SetMode(gconf.GConf.GinMode)

	engine := gin.New()
	engine.Use(gin.Recovery(), middleware.Cors())
	engine.Any("test", response.Request[UserReq, UserRes](UserSrv))
	if err := engine.Run(":8080"); err != nil {
		log.Fatalf("StartHttpServer err: %s", err.Error())
	}
}

// 请求和响应结构体
type UserReq struct {
	server_model.User
	Os   string `json:"os" form:"os" binding:"required" label:"平台"`
	Name string `json:"name" form:"name" binding:"required" label:"姓名"`
}

type UserRes struct {
	Message string `json:"message"`
}

// 服务函数
func UserSrv(ctx context.Context, req *UserReq) (res *UserRes, err error) {
	// 模拟业务逻辑
	//if req.Name == "" {
	//	return nil, fmt.Errorf("name cannot be empty")
	//}
	//initdatabase.GetDefaultClient().
	return &UserRes{Message: fmt.Sprintf("Hello, %s,UserId:%s,Os:%s ", "req.Name", req.UserId, "req.Os")}, nil
}
