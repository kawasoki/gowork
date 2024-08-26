package accerror

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kawasoki/gowork/logger"
	"gorm.io/gorm"
)

//=================================rpcx的错误 =============================================

// Error is customized error.
type Error struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

// Error implements error interface.
func (e *Error) Error() string {
	return fmt.Sprintf(`{"code": %d, "msg": "%s"}`, e.Code, e.Msg)
}

func (e *Error) IsServiceError() bool {
	return true
}

// NewError creates a new Error.
func NewError(code int, msg string) *Error {
	return &Error{
		Code: code,
		Msg:  msg,
	}
}

// NewErrorString creates a new Error from string.
func MewErrorString(s string) (*Error, error) {
	var err Error
	e := json.Unmarshal([]byte(s), &err)

	return &err, e
}

func WrapRpcxError(err error) error {
	if err == nil {
		return nil
	}
	var customErr *Error
	if errors.As(err, &customErr) { //自定义错误 可以返回给客户端的内容,
		err = errors.New(customErr.Msg)
	} else {
		logger.Errorf("rpcx error:%s", err.Error())
		err = ErrInternal
	}
	return err
}

//================================================================

var (
	ErrInternal   = errors.New("internal server error") // 仅供内部使用的错误
	ErrorNotFound = errors.New("error not found")
)

// WrapSystemError 内部错误 不能返回给前端
func WrapSystemError(err error) error {
	if err == nil {
		return nil
	}
	logger.Errorf("server inner error:%s", err.Error())
	return ErrInternal
}

// WrapWarnError 警告日志
func WrapWarnError(err error, msg ...string) error {
	if err == nil {
		return nil
	}
	logger.Warnf("server warn :%s", err.Error())
	if len(msg) == 1 {
		err = errors.New(msg[0])
	}
	return err
	//return fmt.Errorf("%w%w ", ErrWarn, err)
}

//========================================================gormV2的错误====================================

func WrapDbError(err error, ignoreNotFound ...bool) error {
	if err == nil {
		return nil
	}
	//忽略ErrorNotFound信息
	var ignore bool
	if len(ignoreNotFound) == 1 && ignoreNotFound[0] {
		ignore = true
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if ignore {
			return nil
		} else {
			return ErrorNotFound
		}
	}
	logger.Errorf("mysql inner error:%s", err.Error())
	return ErrInternal
}
