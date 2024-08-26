package validator_trans

import (
	"errors"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh2 "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"strings"
)

var trans ut.Translator

func init() {
	translator := zh.New()
	uni := ut.New(translator)
	trans, _ = uni.GetTranslator("zh")
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = zh2.RegisterDefaultTranslations(v, trans)
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("label"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}
}

// 默认开启
var transOff bool

func TurnOff() {
	transOff = true
}

//func InitTrans() {
//	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
//		_ = zh2.RegisterDefaultTranslations(v, trans)
//		v.RegisterTagNameFunc(func(field reflect.StructField) string {
//			name := strings.SplitN(field.Tag.Get("label"), ",", 2)[0]
//			if name == "-" {
//				return ""
//			}
//			return name
//		})
//	}
//}

// Error 多个error只返回第一个error
func Error(err error) string {
	if transOff {
		return "param error"
	}
	var errs validator.ValidationErrors
	if ok := errors.As(err, &errs); !ok {
		return "param error"
	}
	for _, e := range errs {
		errM := e.Translate(trans)
		return errM
	}
	return "param error"
}
