package connection_v2

import (
	"database/sql"
	"fmt"
	mysql2 "github.com/kawasoki/gowork/initdatabase/mysql"
	"github.com/kawasoki/gowork/logger/netlogger"
	mysqlC "github.com/kawasoki/gowork/mysql"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"strings"
	"sync"
	"time"
)

var conns = make(map[string]*gorm.DB)
var mutex sync.Mutex

func GetDbV2(dbConfig *mysql2.DbConfig) *gorm.DB {
	mutex.Lock()
	defer mutex.Unlock()
	tmpDb := conns[dbConfig.DbName]
	if tmpDb == nil {
		//添加详细信息
		name := dbConfig.DbName
		if strings.Contains(dbConfig.DbName, mysqlC.DbProxyConn) {
			name = strings.Split(dbConfig.DbName, mysqlC.DbProxyConn)[1]
		}
		sqlDb, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, name))
		sqlDb.SetMaxIdleConns(dbConfig.MaxIdleConns)
		sqlDb.SetMaxOpenConns(dbConfig.MaxOpenConns)
		sqlDb.SetConnMaxLifetime(time.Second * time.Duration(dbConfig.ConnMaxLife))
		newLogger := logger.New(
			new(GormLogger).Init(),
			//log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             200 * time.Millisecond, // 慢 SQL 阈值
				LogLevel:                  logger.Warn,            // Log level
				IgnoreRecordNotFoundError: true,                   // 忽略ErrRecordNotFound（记录未找到）错误
				Colorful:                  false,                  // 禁用彩色打印
			},
		)
		openedDb, err := gorm.Open(mysql.New(mysql.Config{
			Conn: sqlDb,
		}), &gorm.Config{
			Logger: newLogger.LogMode(logger.Warn),
		})
		if err != nil {
			panic("数据库连接出错：" + err.Error())
		}
		tmpDb = openedDb
		conns[dbConfig.DbName] = tmpDb
	}
	return tmpDb
}

type GormLogger struct {
	logger *zap.Logger
}

func (g *GormLogger) Printf(format string, v ...interface{}) {
	format = strings.ReplaceAll(format, "\n", "")
	if len(v) > 1 {

		if str, ok := v[1].(string); ok {
			if strings.Contains(str, "SLOW") {
				g.logger.Sugar().Warnf(format, v...)
				return
			}
		}
		if _, ok := v[1].(error); ok {
			g.logger.Sugar().Errorf(format, v...)
			return
		}
	}
	g.logger.Sugar().Infof(format, v...)
}
func (g *GormLogger) Init() *GormLogger {
	logAgent := netlogger.ZapLoggerAgent{}
	g.logger = logAgent.Init(&netlogger.LogAgentConf{
		ServerName: "gorm-sql",
		AgentAddr:  "logagent:8899",
		EncoderConf: &zapcore.EncoderConfig{
			MessageKey:       "message",
			LevelKey:         "level",
			EncodeLevel:      zapcore.CapitalLevelEncoder, // INFO
			TimeKey:          "time",
			EncodeTime:       zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
			EncodeCaller:     zapcore.ShortCallerEncoder,
			ConsoleSeparator: " ",
		},
	}).Conn().Daemon().Logger()
	return g
}
