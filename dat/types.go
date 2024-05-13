package dat

import (
	"github.com/892294101/cache-mmap/mmap"
	"github.com/892294101/dds-metadata"
	"github.com/892294101/dds-spfile"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

// WriteCache 写入文件管理
type WriteCache struct {
	ProcName        string // 进程名称
	DatDir          string // 数据目录
	Prefix          string // 文件前缀
	MaxSize         int    // 文件最大size
	dbType          string // 数据库类型
	md              ddsmetadata.MetaData
	pfile           *ddsspfile.Spfile
	file            *mmap.File
	Seq             uint64
	Rba             uint64
	CurrentFile     string
	lock            sync.Mutex
	log             *logrus.Logger
	drain           chan struct{}
	quit            chan struct{}
	wg              sync.WaitGroup
	flushPeriodTime time.Duration
	Dirty           bool
}

// ReadCache 写入文件管理
type ReadCache struct {
	DatDir      string // 数据目录
	Prefix      string // 文件前缀
	pfile       *ddsspfile.Spfile
	file        *mmap.File
	Seq         uint64
	Rba         uint64
	CurrentFile string
}
