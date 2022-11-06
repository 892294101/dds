package connect

import (
	"database/sql"
	"github.com/892294101/dds/spfile"
	oramysql "github.com/892294101/go-mysql/client"
	"github.com/sirupsen/logrus"
)

type Connector interface {
	createConnect() error
	loadEnv() error
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

type ConnectorForOracle struct {
	params  *spfile.Spfile
	log     *logrus.Logger
	dbType  string
	proType string
	conn    *sql.DB
}
