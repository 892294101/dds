package oramysql

import (
	"fmt"
	"github.com/892294101/dds-metadata"
	"github.com/892294101/dds-spfile"
	"github.com/892294101/dds-utils"
	"github.com/892294101/dds/dat"
	"github.com/892294101/dds/ddslog"
	"github.com/892294101/dds/process"
	"github.com/892294101/dds/serialize"
	"github.com/892294101/go-mysql/canal"
	"github.com/892294101/go-mysql/client"
	"github.com/892294101/go-mysql/mysql"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// 同步主体
type Synchronizer struct {
	canalConfig *canal.Config // 同步配置
	DataStream  *canal.Canal  // 同步数据流
	event       *EventExtract
}

type StopService struct {
	power chan int
	force bool
}

type ExtractEvent struct {
	proName           *string                     // 进程名字
	replicationStream *Synchronizer               // 复制流
	pFile             *ddsspfile.Spfile           // 参数文件
	log               *logrus.Logger              // 日志记录器
	client            *client.Conn                // 备用连接
	md                ddsmetadata.MetaData        // 元数据文件
	fileWrite         *dat.WriteCache             // 缓存数据写入
	serialize         map[int]serialize.Serialize // 序列化接口
	rpcServer         *RpcProcess                 // rpc服务器
	stopSrv           *StopService                // 停止进程服务
	isReady           bool                        // 是否准备好检索数据
	rpcPort           int                         //  rpc 端口
	TranHead          *TranHeadFilter             // 数据过滤
}

func (s *Synchronizer) SetSyncHost(host *string, port *uint16) {
	s.canalConfig.Addr = net.JoinHostPort(*host, strconv.Itoa(int(*port)))
}
func (s *Synchronizer) SetSyncUser(user *string)            { s.canalConfig.User = *user }
func (s *Synchronizer) SetSyncPassword(pass *string)        { s.canalConfig.Password = *pass }
func (s *Synchronizer) SetSyncLogger(log *logrus.Logger)    { s.canalConfig.Logger = log }
func (s *Synchronizer) SetSyncFlavor(flavor *string)        { s.canalConfig.Flavor = *flavor }
func (s *Synchronizer) SetServerID(id *uint32)              { s.canalConfig.ServerID = *id }
func (s *Synchronizer) SetRetryConnectNum(nm *int)          { s.canalConfig.MaxReconnectAttempts = *nm }
func (s *Synchronizer) SetCharacter(c *string)              { s.canalConfig.Charset = *c }
func (s *Synchronizer) SetTimeZone(t *time.Location)        { s.canalConfig.TimestampStringLocation = t }
func (s *Synchronizer) SetHeartbeatPeriod(hb time.Duration) { s.canalConfig.HeartbeatPeriod = hb }
func (s *Synchronizer) SetUseDecimal(dm bool)               { s.canalConfig.UseDecimal = dm }
func (s *Synchronizer) SetParseTime(dm bool)                { s.canalConfig.ParseTime = dm }
func (s *Synchronizer) SetSyncAck(se bool)                  { s.canalConfig.SemiSyncEnabled = se }

func (e *ExtractEvent) setProcessName() error {
	pname := e.pFile.GetProcessName()
	if pname == nil {
		return errors.Errorf("set process name failed")
	}
	e.proName = pname
	return nil
}

func (e *ExtractEvent) setSourceDB() error {
	mstr := e.pFile.GetMySQLDBConnStr()
	if mstr != nil {
		// 设置stream连接信息
		e.replicationStream.SetSyncHost(mstr.GetAddress(), mstr.GetPort())
		e.replicationStream.SetSyncUser(mstr.GetUserId())
		e.replicationStream.SetSyncPassword(mstr.GetPassWord())
		e.replicationStream.SetSyncFlavor(mstr.GetTypes())
		e.replicationStream.SetSyncLogger(e.log)
		e.replicationStream.SetServerID(mstr.GetServerID())
		e.replicationStream.SetRetryConnectNum(mstr.GetRetryConnect())
		e.replicationStream.SetCharacter(mstr.GetClientCharacter())
		e.replicationStream.SetTimeZone(mstr.GetTimeZone())
		e.replicationStream.SetHeartbeatPeriod(time.Second * 60)
		e.replicationStream.SetUseDecimal(true)
		e.replicationStream.SetParseTime(true)
		//e.replicationStream.SetSyncAck(true)

		// 初始化备用连接
		conn, err := client.Connect(net.JoinHostPort(*mstr.GetAddress(), strconv.Itoa(int(*mstr.GetPort()))), *mstr.GetUserId(), *mstr.GetPassWord(), "")
		if err != nil {
			return errors.Errorf("standby client connection creation failed: %s", err)
		}

		if err := conn.Ping(); err != nil {
			return errors.Errorf("standby client connection test failed: %s", err)
		}
		e.client = conn
		return nil
	}
	return errors.Errorf("Failed to get %s info", ddsutils.SourceDBType)
}

func (e *ExtractEvent) setLogger(log *logrus.Logger) error {
	if log == nil {
		return errors.Errorf("The logger pointer is null")
	} else {
		e.log = log
		return nil
	}
}

func (e *ExtractEvent) setSpfile(pfile *ddsspfile.Spfile) error {
	if pfile == nil {
		return errors.Errorf("The Spfile pointer is null")
	} else {
		e.pFile = pfile
		return nil
	}
}

// 打开检查点元数据文件
func (e *ExtractEvent) initMetaData(processName string, dataBaseType string, processType string, log *logrus.Logger) error {
	mds, err := ddsmetadata.InitMetaData(processName, dataBaseType, processType, log, ddsmetadata.LOAD)
	if err != nil {
		return err
	}
	e.md = mds
	return nil
}

func (e *ExtractEvent) readPfile(processName string, dataBaseType string, processType string, log *logrus.Logger) error {
	pfile, err := ddsspfile.LoadSpfile(fmt.Sprintf("%s.desc", processName), ddsspfile.UTF8, log, dataBaseType, processType)
	if err != nil {
		return err
	}

	if err := pfile.Production(); err != nil {
		return err
	}
	/*
		// 生成的参数转为json格式，并加载到sqlite数据库，供其它进程调用
		if err := pfile.LoadToDatabase(); err != nil {
			return err
		}
	*/
	if !strings.EqualFold(*pfile.GetProcessName(), processName) {
		return errors.Errorf("Process name mismatch: %s", *pfile.GetProcessName())
	}
	e.pFile = pfile
	return nil
}

func (e *ExtractEvent) GetMasterPosition() (*mysql.Position, error) {
	e.log.Infof("get current file number and position")
	rr, err := e.client.Execute("SHOW MASTER STATUS")
	if err != nil {
		return nil, err
	}
	name, _ := rr.GetString(0, 0)
	pos, _ := rr.GetInt(0, 1)

	e.log.Infof("current file number and position: %v %v", name, pos)
	return &mysql.Position{Name: name, Pos: uint32(pos)}, nil
}

func (e *ExtractEvent) InitSerializeInterface() {
	e.serialize = make(map[int]serialize.Serialize)
	e.serialize[DataEvent] = new(serialize.DataEventV1)
	//e.serialize[DataEvent].InitBuffer()
	e.serialize[TransactionEvent] = new(serialize.TransactionV1)
}

func (e *ExtractEvent) InitSyncerConfig(processName string, dataBaseType string, processType string) {
	processName = strings.ToUpper(processName)
	if len(processName) == 0 {
		e.log.Fatalf("Process name cannot be empty")
	}

	// 初始化同步主体
	e.replicationStream = new(Synchronizer)
	e.replicationStream.canalConfig = new(canal.Config)
	e.replicationStream.event = new(EventExtract)

	// 初始化日志
	log, err := ddslog.InitDDSlog(processName)
	if err != nil {
		fmt.Fprintf(os.Stdout, "%s", err)
		os.Exit(2)
	}

	// 设置抓取进程日志记录器
	if err := e.setLogger(log); err != nil {
		e.log.Fatalf("%s", err)
	}

	// 检查进程文件是否存在
	ok, _ := ddsutils.CheckPcsFile(processName)
	if ok {
		e.log.Fatalf("process group is already running")
	}

	// 读取参数文件
	err = e.readPfile(processName, dataBaseType, processType, log)
	if err != nil {
		e.log.Fatalf("%s", err)
	}

	// 设置进程名称
	if err := e.setProcessName(); err != nil {
		e.log.Fatalf("%s", err)
	}

	// 设置数据库连接地址
	if err := e.setSourceDB(); err != nil {
		e.log.Fatalf("%s", err)
	}

	// 初始化元数据信息
	if err := e.initMetaData(processName, dataBaseType, processType, log); err != nil {
		e.ClearProcessInfoFile()
		e.CloseAll()
		e.log.Fatalf("%s", err)
	}
	// 设置进程启动时间
	if err := e.md.SetStartTime(); err != nil {
		e.ClearProcessInfoFile()
		e.CloseAll()
		e.log.Fatalf("%s", err)
	}

	// 初始化数据流接收队列
	e.replicationStream.event.cache = new(queue)
	if err := e.replicationStream.event.cache.initQueue(log); err != nil {
		e.ClearProcessInfoFile()
		e.CloseAll()
		e.log.Fatalf("%s", err)
	}

	// 文件写入
	fw := dat.NewWriteMgr()
	if err := fw.Init(e.pFile, &dataBaseType, e.md, e.log); err != nil {
		e.ClearProcessInfoFile()
		e.CloseAll()
		e.log.Fatalf("%s", err)
	}

	if err := fw.LoadDatFile(); err != nil {
		e.ClearProcessInfoFile()
		e.CloseAll()
		e.log.Fatalf("%s", err)
	}

	e.fileWrite = fw

	// 初始化序列化接口
	e.InitSerializeInterface()

	// 初始化停止服务
	e.stopSrv = new(StopService)
	e.stopSrv.power = make(chan int, 1)

	// 获取可用的rpc端口
	e.rpcPort, err = ddsutils.GetAvailablePort()
	if err != nil {
		e.ClearProcessInfoFile()
		e.CloseAll()
		e.log.Fatalf("get rpc port error: %s", err)
	}

	// 启动rcp监听
	go e.InitRpc()

	// 启动队列读取
	go e.ProcessCacheData()
}

func (e *ExtractEvent) InitRpc() {
	rs := NewRpc()
	e.rpcServer = rs
	err := rs.StartRpcServer(e, e.log)
	if err != nil {
		e.log.Fatalf("%v", err)
	}
}

func (e *ExtractEvent) StartSyncToStream() {
	/*cfg.IncludeTableRegex = make([]string, 1)
	cfg.IncludeTableRegex[0] = ".*\\.canal_test"
	cfg.ExcludeTableRegex = make([]string, 2)
	cfg.ExcludeTableRegex[0] = "mysql\\..*"
	cfg.ExcludeTableRegex[1] = ".*\\..*_inner"*/

	// fnp := mysql.Position{Name: fmt.Sprintf("mysql-bin.%06d", file), Pos: pos}
	var wait int
ReStart:
	if e.isReady != true {
		wait += 1
		time.Sleep(1 * time.Second)
		if wait >= 7 {
			e.log.Errorf("RPC is incomplete and cannot be start")
			e.ClearProcessInfoFile()
			e.CloseAll()
			e.log.Fatalf("Stop process")
		}
		goto ReStart
	}

	if err := process.WriteProcessInfo(e.pFile, ddsutils.Extract, e.rpcPort); err != nil {
		e.CloseAll()
		e.log.Fatalf("%s", err)
	}

	canStream, err := canal.NewCanal(e.replicationStream.canalConfig)
	if err != nil {
		e.log.Fatalf("Error creating replication stream: %s", err)
	}

	e.replicationStream.DataStream = canStream

	e.replicationStream.DataStream.SetEventHandler(e.replicationStream.event)

	//canStream.RunFrom(mysql.Position{fmt.Sprintf("mysql-bin.%06d", 6),7836403})

	file, pos, err := e.md.GetPosition()
	if err != nil {
		e.log.Fatalf("%s", err)
	}

	if *file == 0 && *pos == 0 {
		pos, err := e.GetMasterPosition()
		if err != nil {
			e.ClearProcessInfoFile()
			e.CloseAll()
			e.log.Fatalf("%s", err)
		}

		s, p, err := ddsutils.ConvertPositionToNumber(pos)
		if err != nil {
			e.ClearProcessInfoFile()
			e.CloseAll()
			e.log.Fatalf("%s", err)
		}

		err = e.md.SetPosition(*s, *p)
		if err != nil {
			e.ClearProcessInfoFile()
			e.CloseAll()
			e.log.Fatalf("%s", err)
		}
		e.log.Infof("Set file number and position for initial startup: %v %v", *s, *p)
		err = e.replicationStream.DataStream.RunFrom(*pos)
		if err != nil {
			e.ClearProcessInfoFile()
			e.CloseAll()
			e.log.Fatalf("%s", err)
		}
	} else if *file != 0 && *pos != 0 {
		position := mysql.Position{Name: fmt.Sprintf("mysql-bin.%06d", *file), Pos: uint32(*pos)}
		err = e.replicationStream.DataStream.RunFrom(position)
		if err != nil {
			e.ClearProcessInfoFile()
			e.CloseAll()
			e.log.Fatalf("%s", err)
		}

	}

}

func (e *ExtractEvent) ClearProcessInfoFile() {
	if err := process.RemoveProcessInfo(e.pFile); err != nil {
		e.log.Errorf("remove process info file failed: %v", err)
	}
}

func (e *ExtractEvent) CloseAll() {
	// 关闭数据写入器（无论是否是否写入完成）
	if e.fileWrite != nil {
		if err := e.fileWrite.CloseDat(); err != nil {
			e.log.Errorf("File writer close error: %v", err)
		}
	}

	// 关闭元数据文件写入器
	if e.md != nil {
		if err := e.md.Close(); err != nil {
			e.log.Errorf("Metadata handle close failed: %v", err)
		}
	}

	// 关闭抓取数据流
	if e.replicationStream != nil && e.replicationStream.DataStream != nil {
		e.replicationStream.DataStream.Close()
	}

	if e.client != nil {
		if err := e.client.Close(); err != nil {
			e.log.Errorf("standby client connection closing failed: %v", err)
		}
	}

	if e.rpcServer != nil {
		e.rpcServer.CloseRpc()
	}
}

func NewMySQLSync() *ExtractEvent {
	return new(ExtractEvent)
}
