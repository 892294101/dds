package process

import (
	"fmt"
	"github.com/892294101/dds/dbs/spfile"
	"github.com/892294101/dds/dbs/utils"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"
)

func WriteProcessInfo(pfile *spfile.Spfile, proType string, rpcPort int) error {
	pid := os.Getpid()
	home, err := utils.GetHomeDirectory()
	if err != nil {
		return err
	}
	pname := pfile.GetProcessName()
	if pname == nil {
		return errors.Errorf("get process name failed the WriteProcessInfo")
	}

	file := filepath.Join(*home, "pcs", *pname)
	ok := utils.IsFileExist(file)
	if ok {
		return errors.Errorf("error writing process group information because process file information already exists")
	}

	hand, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_SYNC, 775)
	if err != nil {
		return errors.Errorf("open file failed the WriteProcessInfo: %v", err)
	}

	var pinfo strings.Builder
	pinfo.WriteString(fmt.Sprintf("%9s: %v\n", utils.PROGRAM, proType))
	pinfo.WriteString(fmt.Sprintf("%9s: %v\n", utils.PROCESSID, *pname))
	pinfo.WriteString(fmt.Sprintf("%9s: %d\n", utils.PORT, rpcPort))
	pinfo.WriteString(fmt.Sprintf("%9s: %d\n", utils.PID, pid))
	pinfo.WriteString(fmt.Sprintf("%9s: %s\n", utils.STATUS, utils.RUNNING))
	pinfo.WriteString(fmt.Sprintf("%9s: %s\n", utils.DBTYPE, spfile.GetMySQLName()))
	hand.WriteString(pinfo.String())
	pinfo.Reset()
	hand.Close()
	return nil
}

func RemoveProcessInfo(pfile *spfile.Spfile) error {
	home, err := utils.GetHomeDirectory()
	if err != nil {
		return err
	}
	pname := pfile.GetProcessName()
	if pname == nil {
		return errors.Errorf("get process name failed the WriteProcessInfo")
	}

	file := filepath.Join(*home, "pcs", *pname)
	ok := utils.IsFileExist(file)
	if ok {
		os.Remove(file)
	}
	return nil
}
