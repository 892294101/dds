package connect

import (
	"database/sql"
	"fmt"
	"github.com/892294101/dds/utils"
	"github.com/godror/godror"
	"github.com/sirupsen/logrus"
	"time"
)

/*func newConnect(p *spfile.Spfile, dbType, proType string, log *logrus.Logger) (*Connector, error) {
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
		case spfile.GetExtractName():
			c = &ConnectorForOracle{params: p, log: log, dbType: dbType, proType: proType}
		default:
			return nil, errors.Errorf("%v process type does not support", proType)
		}
	default:
		return nil, errors.Errorf("%v database type does not support", dbType)
	}

	if err := c.CreateConnect(); err != nil {
		return nil, err
	}

	if err := c.loadEnv(); err != nil {
		return nil, err
	}

	return &c, nil
}*/

/*func (c *ConnectorForMySQL) CreateConnect() error {
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

func (c *ConnectorForMySQL) LoadEnv() error {
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
}*/

func (c *ConnectorForOracle) CreateConnect() (*ConnBody, error) {
	//defer utils.ErrorCheckOfRecover(c.CreateConnect, c.log)

	thread, extractConn, err := c.establish()
	if err != nil {
		return nil, err
	}

	_, standbyConn, err := c.establish()
	if err != nil {
		return nil, err
	}

	c.Conn = &ConnBody{ThreadID: thread, MasterConn: extractConn, AlternativeConn: standbyConn}

	return c.Conn, nil
}

func (c *ConnectorForOracle) establish() (int32, *sql.DB, error) {
	var cs godror.ConnectionParams
	cs.ConnectString = fmt.Sprintf("%v:%d/%v", c.auth.IpAddress, c.auth.Port, c.auth.SID)
	cs.Username = c.auth.UserName
	passwd := godror.NewPassword(c.auth.PassWord)
	cs.Password = passwd
	cs.SessionTimeout = 3
	local, err := time.LoadLocation(c.auth.TimeZone)
	if err != nil {
		return 0, nil, err
	}
	cs.Timezone = local
	cs.SetSessionParamOnInit("NLS_NUMERIC_CHARACTERS", ",.")

	language, territory, charset, err := utils.ParseNLSLANG(c.auth.Character)
	if err != nil {
		return 0, nil, err
	} else {
		cs.SetSessionParamOnInit("NLS_LANGUAGE", language)
		cs.SetSessionParamOnInit("NLS_TERRITORY", territory)
		cs.SetSessionParamOnInit("NLS_CHARACTERSET", charset)
	}

	cs.SetSessionParamOnInit("NLS_DATE_FORMAT", "YYYY-MM-DD HH24:MI:SS")
	cs.SetSessionParamOnInit("NLS_TIMESTAMP_FORMAT", "YYYY-MM-DD HH24:MI:SSXFF")
	cs.SetSessionParamOnInit("NLS_TIMESTAMP_TZ_FORMAT", "YYYY-MM-DD HH24:MI:SSXFF TZR")
	// db.Exec("ALTER SESSION SET OCI_SESSION_STATELESS = TRUE")
	conn := sql.OpenDB(godror.NewConnector(cs))
	conn.SetMaxIdleConns(MaxConnectionConcurrent)
	conn.SetMaxOpenConns(MaxConnectionConcurrent)

	if err := conn.Ping(); err != nil {
		return 0, nil, err
	}

	if err := c.loadEnv(); err != nil {
		return 0, nil, err
	}

	res := conn.QueryRow("select thread# from v$instance")
	var thread int32
	if err := res.Scan(&thread); err != nil {
		return 0, nil, err
	}

	return thread, conn, nil
}

func (c *ConnectorForOracle) SetAuth(a *Auth, l *logrus.Logger) {
	c.auth = a
	c.log = l
}

func (c *ConnectorForOracle) loadEnv() error {
	/*	SQL> select userenv('language') from dual;
		USERENV('LANGUAGE')
		----------------------------------------------------
		AMERICAN_AMERICA.ZHS16GBK
	*/

	return nil
}

// Close 它会关闭主连接和备用连接
func (c *ConnBody) Close() error {
	var merr error
	var serr error
	if c.MasterConn != nil {
		merr = c.MasterConn.Close()
	} else {
		merr = fmt.Errorf("master connect null pointer")
	}

	if c.AlternativeConn != nil {
		serr = c.MasterConn.Close()
	} else {
		serr = fmt.Errorf("alternative connect null pointer")
	}
	if merr != nil || serr != nil {
		return fmt.Errorf("%s:%s", merr, serr)
	}
	return nil
}
