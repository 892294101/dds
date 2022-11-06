package connect

import (
	"database/sql"
	"fmt"
	"github.com/892294101/dds/spfile"
	"github.com/892294101/dds/utils"
	oramysql "github.com/892294101/go-mysql/client"
	"github.com/godror/godror"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"time"
)

func newConnect(p *spfile.Spfile, dbType, proType string, log *logrus.Logger) (*Connector, error) {
	if p == nil {
		return nil, errors.Errorf("spfile file is not empty")
	}
	if log == nil {
		return nil, errors.Errorf("Logger file is not empty")
	}

	if len(dbType) == 0 {
		return nil, errors.Errorf("Database type cannot be empty")
	}

	if len(proType) == 0 {
		return nil, errors.Errorf("Process type cannot be empty")
	}

	var c Connector
	switch dbType {
	case spfile.GetMySQLName():
		switch proType {
		case spfile.GetReplicationName():
			c = &ConnectorForMySQL{params: p, log: log, dbType: dbType, proType: proType}
		default:
			return nil, errors.Errorf("%v process type does not support", proType)
		}
	case spfile.GetOracleName():
		switch proType {
		case spfile.GetReplicationName():
			if len(p.GetOracleDBConnStr().GetHostAddress()) != 1 {
				return nil, errors.Errorf("Oracle replication process can only configure one target IP address")
			}
			c = &ConnectorForOracle{params: p, log: log, dbType: dbType, proType: proType}
		default:
			return nil, errors.Errorf("%v process type does not support", proType)
		}
	default:
		return nil, errors.Errorf("%v database type does not support", dbType)
	}

	if err := c.createConnect(); err != nil {
		return nil, err
	}

	if err := c.loadEnv(); err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *ConnectorForMySQL) createConnect() error {
	defer utils.ErrorCheckOfRecover(c.createConnect, c.log)
	c.log.Infof("create database connection for host address: %v", *c.params.GetMySQLDBConnStr().GetAddress())

	conn, err := oramysql.Connect(
		*c.params.GetMySQLDBConnStr().GetAddress(),
		*c.params.GetMySQLDBConnStr().GetUserId(),
		*c.params.GetMySQLDBConnStr().GetPassWord(),
		"",
	)
	if err != nil {
		return err
	}

	if err := conn.Ping(); err != nil {
		return err
	}
	c.conn = conn

	if err := c.loadEnv(); err != nil {
		return err
	}

	return nil
}

func (c *ConnectorForMySQL) loadEnv() error {
	c.log.Infof("loading database environment")
	_, err := c.conn.Execute("SET AUTOCOMMIT = 0")
	if err != nil {
		return err
	}

	err = c.conn.SetCharset(*c.params.GetMySQLDBConnStr().GetClientCharacter())
	if err != nil {
		return err
	}

	_, err = c.conn.Execute(fmt.Sprintf("SET COLLATION_CONNECTION = %s", *c.params.GetMySQLDBConnStr().GetClientCollation()))
	if err != nil {
		return err
	}
	return nil
}

func (c *ConnectorForOracle) createConnect() error {
	defer utils.ErrorCheckOfRecover(c.createConnect, c.log)
	c.log.Infof("create database connection for host address: %v", *c.params.GetMySQLDBConnStr().GetAddress())

	var cs godror.ConnectionParams
	cs.ConnectString = fmt.Sprintf("%v:%d/%v", c.params.GetOracleDBConnStr().GetHostAddress()[0], c.params.GetOracleDBConnStr().GetPort(), c.params.GetOracleDBConnStr().GetSID())
	cs.Username = c.params.GetOracleDBConnStr().GetUserName()
	passwd := godror.NewPassword(c.params.GetOracleDBConnStr().GetPassWord())
	cs.Password = passwd
	cs.SessionTimeout = 3
	local, err := time.LoadLocation(c.params.GetOracleDBConnStr().GetTimeZone())
	if err != nil {
		return err
	}
	cs.Timezone = local
	cs.SetSessionParamOnInit("NLS_NUMERIC_CHARACTERS", ",.")
	cs.SetSessionParamOnInit("NLS_LANGUAGE", "AMERICAN")
	cs.SetSessionParamOnInit("NLS_TERRITORY", "AMERICA")
	cs.SetSessionParamOnInit("NLS_CHARACTERSET", "ZHS16GBK")
	cs.SetSessionParamOnInit("NLS_DATE_FORMAT", "YYYY-MM-DD HH24:MI:SS")
	cs.SetSessionParamOnInit("NLS_TIMESTAMP_FORMAT", "YYYY-MM-DD HH24:MI:SSXFF")
	cs.SetSessionParamOnInit("NLS_TIMESTAMP_TZ_FORMAT", "YYYY-MM-DD HH24:MI:SSXFF TZR")

	conn := sql.OpenDB(godror.NewConnector(cs))
	conn.SetMaxIdleConns(MaxConnectionConcurrent)
	conn.SetMaxOpenConns(MaxConnectionConcurrent)

	if err := conn.Ping(); err != nil {
		return err
	}
	c.conn = conn

	if err := c.loadEnv(); err != nil {
		return err
	}

	return nil
}

func (c *ConnectorForOracle) loadEnv() error {
	/*	SQL> select userenv('language') from dual;
		USERENV('LANGUAGE')
		----------------------------------------------------
		AMERICAN_AMERICA.ZHS16GBK
	*/

	return nil
}
