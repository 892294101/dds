package connect

import (
	"database/sql"
	"github.com/892294101/dds/spfile"
	oramysql "github.com/892294101/go-mysql/client"
	"github.com/sirupsen/logrus"
)

type Connector interface {
	CreateConnect() (*ConnBody, error)
	SetAuth(a *Auth, l *logrus.Logger)
}

const (
	MaxConnectionConcurrent = 1
)

type ConnectorForMySQL struct {
	params  *spfile.Spfile
	log     *logrus.Logger
	dbType  string
	proType string
	conn    *oramysql.Conn
}

type ConnBody struct {
	ThreadID        int32
	MasterConn      *sql.DB
	AlternativeConn *sql.DB
}

type ConnectorForOracle struct {
	log  *logrus.Logger // 日志系统
	Conn *ConnBody      // 连接
	auth *Auth          // 认证信息
}

type Auth struct {
	IpAddress string // 主机地址
	UserName  string // 用户名
	PassWord  string // 密码
	SID       string // SID
	Port      uint16 // 端口
	Character string // 字符集
	Retry     int    // 连接重试次数
	TimeZone  string // 时区
}
