package oracle

import (
	"github.com/892294101/dds/dbs/connect"
	"github.com/892294101/dds/dbs/spfile"
	"github.com/sirupsen/logrus"
)

// CaptureTasks 捕获任务
type CaptureTasks struct {
	extInstanceBody map[int]*ExtInstanceBody // 实例Key编码
	paramfile       *spfile.Spfile           // 参数文件
	log             *logrus.Logger           // 日志记录器
}

// ExtInstanceBody 抓取实例体，默认最高4个线程
type ExtInstanceBody struct {
	extThread []*ExtThread  // 每个实例抓取的线程实体
	auth      *connect.Auth // 当前实例认证
	threadId  int32         // 当前实例日志线程id
	sequence  int32         // 当前实例日志序列
	firstSCN  uint64        // 当前实例日志开始SCN
	endSCN    uint64        // 当前实例日志结束SCN
	logfile   string        // 当前实例日志文件
}

// ExtThread 抓取线程实体
type ExtThread struct {
	eid      int                 // 抓取线程id编码（app内）
	db       *connect.ConnBody   // 数据库连接
	area     chan *logContentV11 // 缓存区
	state    int8                // 抓取状态：0：从未使用，1：正在解析中，2：暂停中，3：空闲
	threadId int32               // 日志文件所属线程
	sequence int32               // 日志文件序列号
	logfile  string              // 日志文件（归档或者redo）
	firstSCN uint64              // 文件开始scn号
	endSCN   uint64              // 文件结束scn号
}

// 初始化抓取实例
func (t *CaptureTasks) newExtInstanceBody() error {
	// 抓取线程默认值
	Parallel := 4

	inc := len(t.paramfile.GetOracleDBConnStr().GetHostAddress())
	t.extInstanceBody = make(map[int]*ExtInstanceBody, inc) // 根据IP数量决定数据库实例数
	// 初始化抓取线程
	for i := 0; i < inc; i++ {
		t.extInstanceBody[i+1] = &ExtInstanceBody{
			extThread: make([]*ExtThread, Parallel),
			auth: &connect.Auth{
				IpAddress: t.paramfile.GetOracleDBConnStr().GetHostAddress()[i],
				UserName:  t.paramfile.GetOracleDBConnStr().GetUserName(),
				PassWord:  t.paramfile.GetOracleDBConnStr().GetPassWord(),
				SID:       t.paramfile.GetOracleDBConnStr().GetSID(),
				Port:      t.paramfile.GetOracleDBConnStr().GetPort(),
				Character: t.paramfile.GetOracleDBConnStr().GetClientCharacter(),
				Retry:     t.paramfile.GetOracleDBConnStr().GetRetryMaxConnNumber(),
				TimeZone:  t.paramfile.GetOracleDBConnStr().GetTimeZone(),
			},
		}
	}

	return t.reconnectDB()
}

// fixedLogFile 定位日志文件
func (t *CaptureTasks) fixedLogFile() error {
	return nil
}

// reconnectDB 连接数据库
func (t *CaptureTasks) reconnectDB() error {
	for _, body := range t.extInstanceBody {
		for i, thread := range body.extThread {
			var cn connect.Connector
			cn = &connect.ConnectorForOracle{}
			cn.SetAuth(body.auth, t.log)
			conn, err := cn.CreateConnect()
			if err != nil {
				return err
			}
			thread.threadId = int32(i)
			thread.db = conn
		}
	}

	return nil
}

// InitTaskGroups 初始化抓取任务
func (t *CaptureTasks) InitTaskGroups(s *spfile.Spfile, log *logrus.Logger) error {
	t.log = log
	t.paramfile = s

	return t.newExtInstanceBody()
}

// NewTaskGroups new 抓取任务
func NewTaskGroups() *CaptureTasks {
	return new(CaptureTasks)
}
