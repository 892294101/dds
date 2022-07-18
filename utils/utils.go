package utils

import (
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// 支持的数据库和进程类型
const (
	MySQL    = "MySQL"
	MariaDB  = "MariaDB"
	Oracle   = "Oracle"
	Extract  = "Extract"
	Replicat = "Replicat "
)

// process参数
var (
	ProcessType    = "PROCESS" // 参数类型
	ProcessRegular = "(^)(?i:(" + ProcessType + "))(\\s+)((?:[A-Za-z0-9_]){4,12})($)"
)

// sourcedb参数
var (
	SourceDBType = "SOURCEDB"  // 参数类型
	Port         = "PORT"      // 端口关键字
	DataBase     = "DATABASE"  // 默认连接的数据库
	Types        = "TYPE"      // 库类型,可选mysql mariadb
	UserId       = "USERID"    // 连接用户
	PassWord     = "PASSWORD"  // 连接密码
	ServerId     = "SERVERID"  // mysql server id
	Retry        = "RETRY"     // 连接重试最大
	Character    = "CHARACTER" // 客户端字符集关键字
	Collation    = "COLLATION"

	DefaultPort            uint16 = 3306    // 默认端口
	DefaultDataBase               = "test"  // 默认连接数据库
	DefaultTypes                  = "mysql" // 默认库类型
	DefaultUserId                 = "root"  // 默认用户名
	DefaultMaxRetryConnect        = 3
	DefaultClientCharacter        = "UTF8"
	DefaultClientCollation        = "UTF8_GENERAL_CI"
)

// traildir 参数
var (
	TrailDirType          = "TRAILDIR" // 参数类型
	TrailSizeKey          = "SIZE"     // size关键字
	TrailKeepKey          = "KEEP"     // keey 关键字
	MB                    = "MB"
	GB                    = "GB"
	DAY                   = "DAY"
	DefaultTrailMaxSize   = 128 // 默认trail文件的最大, 单位是M, 单位不可更改
	DefaultTrailMinSize   = 16  // 默认trail文件的最小
	DefaultTrailKeepValue = 7   // 默认trail文件保留时间,默认是天
)

// discardfile 参数
var (
	DiscardFileType    = "DISCARDFILE"
	DiscardFileRegular = "(^)(?i:(" + DiscardFileType + "))(\\s+)((.+))($)"
)

// dboptions 参数
var (
	DBOptionsType      = "DBOPTIONS"
	SuppressionTrigger = "SUPPRESSIONTRIGGER" // 表操作时抑制触发器
	IgnoreReplicates   = "IGNOREREPLICATES"   // 忽略复制进程执行的操作
	GetReplicates      = "GETREPLICATES"      // 获取复制进程的操作
	IgnoreForeignkey   = "IGNOREFOREIGNKEY"   // 忽略外键约束
)

// TABLE 参数
var (
	TableType    = "TABLE"
	TableRegular = "(^)(?i:(" + TableType + "))(\\s+)((\\S+)(\\.)(\\S+\\s*)(;))($)"
)

// TABLEExclude 参数
var (
	TableExcludeType    = "TABLEEXCLUDE"
	TableExcludeRegular = "(^)(?i:(" + TableExcludeType + "))(\\s+)((\\S+)(\\.)(\\S+\\s*)(;))($)"
)

//根据执行文件路径获取程序的HOME路径
func GetHomeDirectory() (s *string, err error) {
	file, _ := exec.LookPath(os.Args[0])
	ExecFilePath, _ := filepath.Abs(file)
	var dir string

	os := runtime.GOOS
	switch os {
	case "windows":
		execfileslice := strings.Split(ExecFilePath, `\`)
		HomeDirectory := execfileslice[:len(execfileslice)-2]
		for i, v := range HomeDirectory {
			if v != "" {
				if i > 0 {
					dir += `\` + v
				} else {
					dir += v
				}
			}
		}
	case "linux":
		execfileslice := strings.Split(ExecFilePath, "/")
		HomeDirectory := execfileslice[:len(execfileslice)-2]
		for _, v := range HomeDirectory {
			if v != "" {
				dir += `/` + v
			}
		}
	default:
		return nil, errors.Errorf("Unsupported operating system type: %s", os)
	}

	if dir == "" {
		return nil, errors.Errorf("Get program home directory failed: %s", dir)
	}
	return &dir, nil
}

func HasPrefixIgnoreCase(s, prefix string) bool {
	return len(s) >= len(prefix) && strings.EqualFold(s[0:len(prefix)], prefix)
}

func TrimKeySpace(s []string) []string {
	var deDup []string
	for _, rv := range s {
		if strings.TrimSpace(rv) != "" {
			deDup = append(deDup, strings.TrimSpace(rv))
		}
	}
	return deDup
}

func KeyCheck(s *string) bool {
	key := map[string]string{
		strings.ToUpper(SourceDBType): SourceDBType,
		strings.ToUpper(Port):         Port,
		strings.ToUpper(DataBase):     DataBase,
		strings.ToUpper(Types):        Types,
		strings.ToUpper(UserId):       UserId,
		strings.ToUpper(PassWord):     PassWord,
		strings.ToUpper(TrailDirType): TrailDirType,
		strings.ToUpper(TrailSizeKey): TrailSizeKey,
		strings.ToUpper(TrailKeepKey): TrailKeepKey,
	}
	_, ok := key[strings.ToUpper(*s)]
	return ok
}
