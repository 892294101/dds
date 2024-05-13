package process

import (
	"fmt"
	"github.com/892294101/dds-spfile"
	"github.com/892294101/dds-utils"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"
)

func WriteProcessInfo(pfile *ddsspfile.Spfile, proType string, rpcPort int) error {
	pid := os.Getpid()
	home, err := ddsutils.GetHomeDirectory()
	if err != nil {
		return err
	}
	pname := pfile.GetProcessName()
	if pname == nil {
		return errors.Errorf("get process name failed the WriteProcessInfo")
	}

	file := filepath.Join(*home, "pcs", *pname)
	ok := ddsutils.IsFileExist(file)
	if ok {
		return errors.Errorf("error writing process group information because process file information already exists")
	}

	hand, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_SYNC, 775)
	if err != nil {
		return errors.Errorf("open file failed the WriteProcessInfo: %v", err)
	}

	var pinfo strings.Builder
	pinfo.WriteString(fmt.Sprintf("%9s: %v\n", ddsutils.PROGRAM, proType))
	pinfo.WriteString(fmt.Sprintf("%9s: %v\n", ddsutils.PROCESSID, *pname))
	pinfo.WriteString(fmt.Sprintf("%9s: %d\n", ddsutils.PORT, rpcPort))
	pinfo.WriteString(fmt.Sprintf("%9s: %d\n", ddsutils.PID, pid))
	pinfo.WriteString(fmt.Sprintf("%9s: %s\n", ddsutils.STATUS, ddsutils.RUNNING))
	pinfo.WriteString(fmt.Sprintf("%9s: %s\n", ddsutils.DBTYPE, ddsspfile.GetMySQLName()))
	hand.WriteString(pinfo.String())
	pinfo.Reset()
	hand.Close()
	return nil
}

func RemoveProcessInfo(pfile *ddsspfile.Spfile) error {
	home, err := ddsutils.GetHomeDirectory()
	if err != nil {
		return err
	}
	pname := pfile.GetProcessName()
	if pname == nil {
		return errors.Errorf("get process name failed the WriteProcessInfo")
	}

	file := filepath.Join(*home, "pcs", *pname)
	ok := ddsutils.IsFileExist(file)
	if ok {
		os.Remove(file)
	}
	return nil
}
