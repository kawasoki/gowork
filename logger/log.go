package logger

import (
	"github.com/kawasoki/gowork/logger/netlogger"
	"go.uber.org/zap"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var logger *zap.Logger
var once sync.Once
var ServerName = "mobile-app"

func init() {
	once.Do(func() {
		executablePath, err := os.Executable()
		if err != nil {
			log.Fatalf("初始化日志logger失败:%s", err.Error())
		}
		if name := filepath.Base(executablePath); name != "" {
			ServerName = name
		} else {
			log.Println("没有获取到服务名称")
		}
		InitLog(ServerName, "logagent:8899")
	})
}

func InitLog(serverName string, addr string) {
	logConf := &netlogger.LogAgentConf{
		ServerName: serverName,
		AgentAddr:  addr,
	}
	logger = new(netlogger.ZapLoggerAgent).Init(logConf).Conn().Daemon().Logger()
}

func Error(args ...interface{}) {
	logger.Sugar().Error(args...)
}
func Errorf(template string, args ...interface{}) {
	logger.Sugar().Errorf(template, args...)
}
func Info(args ...interface{}) {
	logger.Sugar().Info(args...)
}
func Infof(template string, args ...interface{}) {
	logger.Sugar().Infof(template, args...)
}
func Warn(args ...interface{}) {
	logger.Sugar().Warn(args...)
}
func Warnf(template string, args ...interface{}) {
	logger.Sugar().Warnf(template, args...)
}
