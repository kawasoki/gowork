package mysql

type DbConfig struct {
	Username     string
	Password     string
	Host         string
	Port         string
	DbName       string
	MaxIdleConns int
	MaxOpenConns int
	SqlOutWirte  int
	ConnMaxLife  int
}

// 只能通过这种方式获取配置对象
func NewDbConfig() *DbConfig {
	conf := &DbConfig{}
	return conf
}
