package utils

import (
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
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
