package dbv2

import (
	"github.com/kawasoki/gowork/gconf"
	"github.com/kawasoki/gowork/initdatabase/mysql"
	"github.com/kawasoki/gowork/initdatabase/mysql/connection_v2"
	"gorm.io/gorm"
)

var (
	g_UserDbConf *mysql.DbConfig
)

// 得到当前数据库连接
func GetUserDb() *gorm.DB {

	if g_UserDbConf == nil {
		g_UserDbConf = mysql.NewDbConfig()
		g_UserDbConf.DbName = gconf.GConf.UserDbName
		g_UserDbConf.Host = gconf.GConf.MysqlHost
		g_UserDbConf.Username = gconf.GConf.MysqlUser
		g_UserDbConf.Password = gconf.GConf.MysqlPwd
		g_UserDbConf.Port = gconf.GConf.MysqlPort
		g_UserDbConf.MaxIdleConns = gconf.GConf.MysqlMaxIdleConns
		g_UserDbConf.MaxOpenConns = gconf.GConf.MysqlMaxOpenConns
		g_UserDbConf.SqlOutWirte = gconf.GConf.SqlOutWirte
		g_UserDbConf.ConnMaxLife = gconf.GConf.MysqlConnMaxLife
	}

	return connection_v2.GetDbV2(g_UserDbConf)
}
