package oramysql

import (
	"context"
	"fmt"
	cache "github.com/emirpasic/gods/queues/linkedlistqueue"
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"myGithubLib/dds/extract/mysql/spfile"
	"myGithubLib/dds/extract/mysql/utils"
	"time"
)

const (
	SPfileUninitError = "SPfile Parameter not initialized"
)

/*
修改源码 func (s *BinlogStreamer) closeWithError(err error) { 删除非日志接口的错误 78行
修改logger 日志记录器
*/

// 同步主体
type Synchronizer struct {
	syncerConfig   *replication.BinlogSyncerConfig // 同步配置
	binlogSyncer   *replication.BinlogSyncer       // 同步器
	binlogStreamer *replication.BinlogStreamer     // 同步流
}

type ExtractEvent struct {
	proName *string        // 进程名字
	syncBus *Synchronizer  // 复制流
	cache   *cache.Queue   // 缓存队列
	pFile   *spfile.Spfile // 参数文件
	log     *logrus.Logger // 日志记录器
}

func (s *Synchronizer) SetSyncHost(host *string)         { s.syncerConfig.Host = *host }
func (s *Synchronizer) SetSyncPort(port *uint16)         { s.syncerConfig.Port = *port }
func (s *Synchronizer) SetSyncUser(user *string)         { s.syncerConfig.User = *user }
func (s *Synchronizer) SetSyncPassword(pass *string)     { s.syncerConfig.Password = *pass }
func (s *Synchronizer) SetSyncLogger(log *logrus.Logger) { s.syncerConfig.Logger = log }
func (s *Synchronizer) SetSyncFlavor(flavor *string)     { s.syncerConfig.Flavor = *flavor }
func (s *Synchronizer) SetServerID(id *uint32)           { s.syncerConfig.ServerID = *id }
func (s *Synchronizer) SetRetryConnectNum(nm *int)       { s.syncerConfig.MaxReconnectAttempts = *nm }
func (s *Synchronizer) SetCharacter(c *string)           { s.syncerConfig.Charset = *c }

func (s *Synchronizer) NewBinlogSyncer() {
	s.binlogSyncer = replication.NewBinlogSyncer(*s.syncerConfig)
}

func (s *Synchronizer) StartSync(fnp mysql.Position) error {
	streams, err := s.binlogSyncer.StartSync(fnp)
	if err != nil {
		return err
	}
	s.binlogStreamer = streams
	return nil
}

func (e *ExtractEvent) setProcessName() error {
	if e.pFile.ParamSet == nil {
		return errors.Errorf(SPfileUninitError)
	}
	res := e.pFile.ParamSet[utils.ProcessType].GetParam()
	if res != nil {
		switch v := res.(type) {
		case *spfile.Process:
			e.proName = v.ProInfo.GetName()
			return nil
		}
	}
	return errors.Errorf("Failed to get %s info", utils.SourceDBType)
}

func (e *ExtractEvent) setSourceDB() error {
	if e.pFile.ParamSet == nil {
		return errors.Errorf(SPfileUninitError)
	}
	res := e.pFile.ParamSet[utils.SourceDBType].GetParam()
	if res != nil {
		switch v := res.(type) {
		case *spfile.SourceDB:
			e.syncBus.SetSyncHost(v.DBInfo.GetAddress())
			e.syncBus.SetSyncPort(v.DBInfo.GetPort())
			e.syncBus.SetSyncUser(v.DBInfo.GetUserId())
			e.syncBus.SetSyncPassword(v.DBInfo.GetPassWord())
			e.syncBus.SetSyncFlavor(v.DBInfo.GetTypes())
			e.syncBus.SetSyncLogger(e.log)
			e.syncBus.SetServerID(v.DBInfo.GetServerID())
			e.syncBus.SetRetryConnectNum(v.DBInfo.GetRetryConnect())
			e.syncBus.SetCharacter(v.DBInfo.GetClientCharacter())
		}
		return nil
	}
	return errors.Errorf("Failed to get %s info", utils.SourceDBType)
}

func (e *ExtractEvent) setLogger(log *logrus.Logger) error {
	if log == nil {
		return errors.Errorf("The logger pointer is null")
	} else {
		e.log = log
		return nil
	}
}

func (e *ExtractEvent) setSpfile(pfile *spfile.Spfile) error {
	if pfile == nil {
		return errors.Errorf("The Spfile pointer is null")
	} else {
		e.pFile = pfile
		return nil
	}
}

func (e *ExtractEvent) InitSyncerConfig(log *logrus.Logger, pfile *spfile.Spfile) error {
	e.syncBus = new(Synchronizer)                                // 初始化同步主体
	e.syncBus.syncerConfig = new(replication.BinlogSyncerConfig) // 初始化

	// 设置日志记录器
	if err := e.setLogger(log); err != nil {
		return err
	}

	// 设置spfile参数
	if err := e.setSpfile(pfile); err != nil {
		return err
	}

	// 设置进程名称
	if err := e.setProcessName(); err != nil {
		return err
	}

	// 设置数据库连接地址
	if err := e.setSourceDB(); err != nil {
		return err
	}

	// 创建一个binlog同步
	e.syncBus.NewBinlogSyncer()
	return nil
}

func (e *ExtractEvent) StartSyncToStream(file int, pos uint32) error {
	fnp := mysql.Position{Name: fmt.Sprintf("mysql-bin.%06d", file), Pos: pos}

	if err := e.syncBus.StartSync(fnp); err != nil {
		return err
	}
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		ev, err := e.syncBus.binlogStreamer.GetEvent(ctx)
		cancel()

		if err == context.DeadlineExceeded {
			continue
		} else if err == replication.ErrNeedSyncAgain || err == replication.ErrSyncClosed {
			e.syncBus.binlogSyncer.Close()
			break
		} else {
			if err != nil {
				e.log.Fatalf("%s", err)
			}
		}
		/*if ev != nil {
			ev.Dump(os.Stdout)
		}*/
		switch v := ev.Event.(type) {
		case *replication.QueryEvent:
			fmt.Println(string(v.Query) )
		}

	}

	return nil
}
func NewMySQLSync() *ExtractEvent {
	return new(ExtractEvent)
}

func NewCanalConfig() *canal.Config {

	return nil
}
