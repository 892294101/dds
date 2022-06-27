package utils

import (
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	MySQL    = "MySQL"
	MariaDB  = "MariaDB"
	Oracle   = "Oracle"
	Extract  = "Extract"
	Replicat = "Replicat "
)

var (
	ProcessType    = "PROCESS"
	ProcessRegular = "(^)(?i:(" + ProcessType + "))(\\s+)((?:[A-Za-z0-9_]){4,12})($)"
)

var (
	SourceDBType = "SOURCEDB"
	Port         = "PORT"
	DataBase     = "DATABASE"
	Types        = "TYPE"
	UserId       = "USERID"
	PassWord     = "PASSWORD"
)

var (
	TrailDirType = "TRAILDIR"
	//TrailDirRegular = "(^)(?i:(" + TrailDirType + "))(\\s+)((.+))($)"
	TrailSizeKey          = "SIZE"
	TrailKeepKey          = "KEEP"
	MB                    = "MB"
	GB                    = "GB"
	DAY                   = "DAY"
	DefaultTrailSize      = 128
	DefaultTrailKeepValue = 7
)

var (
	DefaultPort     = "3306"
	DefaultDataBase = "test"
	DefaultTypes    = "mysql"
	DefaultUserId   = "root"
)

//根据执行文件路径获取程序的HOME路径
func GetHomeDirectory() (dir *string, err error) {
	file, _ := exec.LookPath(os.Args[0])
	ExecFilePath, _ := filepath.Abs(file)

	os := runtime.GOOS
	switch os {
	case "windows":
		execfileslice := strings.Split(ExecFilePath, `\`)
		HomeDirectory := execfileslice[:len(execfileslice)-2]
		for i, v := range HomeDirectory {
			if v != "" {
				if i > 0 {
					*dir += `\` + v
				} else {
					*dir += v
				}
			}
		}
	case "linux":
		execfileslice := strings.Split(ExecFilePath, "/")
		HomeDirectory := execfileslice[:len(execfileslice)-2]
		for _, v := range HomeDirectory {
			if v != "" {
				*dir += `/` + v
			}
		}
	default:
		return nil, errors.Errorf("Unsupported operating system type: %s", os)
	}

	if *dir == "" {
		return nil, errors.Errorf("Get program home directory failed: %s", dir)
	}
	return dir, nil
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
