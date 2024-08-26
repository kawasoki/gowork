package response

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/kawasoki/gowork/configs"
	"github.com/kawasoki/gowork/logger"
	"io"
	"net/http"
)

// 泛型服务函数类型
type SrvFunc[SrvReq any, SrvRes any] func(ctx context.Context, req SrvReq) (res SrvRes, err error)

// Request 处理请求
func Request[Req any, Res any](serviceFunc SrvFunc[*Req, *Res]) gin.HandlerFunc {
	return func(c *gin.Context) {
		response := Response{Code: http.StatusOK, Msg: http.StatusText(http.StatusOK)}
		ctx, cancel := context.WithTimeout(context.Background(), configs.DefaultHttpTimeout)
		defer cancel()
		req := new(Req)
		if _, ok := any(req).(EmptyMarker); !ok { //请求参数不为空if !reflect.DeepEqual(req, &Empty{}) {
			_ = c.ShouldBindHeader(req) //在下一步处理请求参数的时候验证 binding
			//params
			if err := c.ShouldBind(req); err != nil && !errors.Is(err, io.EOF) {
				reqRaw, _ := json.Marshal(req)
				logger.Warnf("params invalid, err: %s, reqRaw: %s", err.Error(), reqRaw)
				c.JSON(http.StatusOK, response.handleValidateErr(err))
				return
			}
			//token解析出来的userid
			if userId := c.GetString(common.UserId); userId != "" {
				//这种方式 定义req结构体的时候 user_id用 server_model.User
				if u, o := any(req).(interface{ SetUserId(string) }); o {
					u.SetUserId(userId)
				}
				//避免使用反射
				//reqValue := reflect.ValueOf(req).Elem()
				//if field := reqValue.FieldByName("UserId"); field.IsValid() && field.CanSet() && field.Kind() == reflect.String {
				//	field.SetString(userId)
				//}
			}
			//处理 page pageSize
			if p, o1 := any(req).(interface{ CheckPage() }); o1 {
				p.CheckPage()
			}
		}
		if res, err := serviceFunc(ctx, req); err != nil { //调用服务函数
			response.Code = http.StatusBadRequest
			response.Msg = err.Error()
		} else {
			if res != nil { //reflect.ValueOf(res).IsZero() 结构体零值也为true
				response.Data = res
			}
		}
		c.JSON(http.StatusOK, response)
	}
}
