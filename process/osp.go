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

func WriteProcessInfo(pfile *dds_spfile.Spfile, proType string, rpcPort int) error {
	pid := os.Getpid()
	home, err := dds_utils.GetHomeDirectory()
	if err != nil {
		return err
	}
	pname := pfile.GetProcessName()
	if pname == nil {
		return errors.Errorf("get process name failed the WriteProcessInfo")
	}

	file := filepath.Join(*home, "pcs", *pname)
	ok := dds_utils.IsFileExist(file)
	if ok {
		return errors.Errorf("error writing process group information because process file information already exists")
	}

	hand, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_SYNC, 775)
	if err != nil {
		return errors.Errorf("open file failed the WriteProcessInfo: %v", err)
	}

	var pinfo strings.Builder
	pinfo.WriteString(fmt.Sprintf("%9s: %v\n", dds_utils.PROGRAM, proType))
	pinfo.WriteString(fmt.Sprintf("%9s: %v\n", dds_utils.PROCESSID, *pname))
	pinfo.WriteString(fmt.Sprintf("%9s: %d\n", dds_utils.PORT, rpcPort))
	pinfo.WriteString(fmt.Sprintf("%9s: %d\n", dds_utils.PID, pid))
	pinfo.WriteString(fmt.Sprintf("%9s: %s\n", dds_utils.STATUS, dds_utils.RUNNING))
	pinfo.WriteString(fmt.Sprintf("%9s: %s\n", dds_utils.DBTYPE, dds_spfile.GetMySQLName()))
	hand.WriteString(pinfo.String())
	pinfo.Reset()
	hand.Close()
	return nil
}

func RemoveProcessInfo(pfile *dds_spfile.Spfile) error {
	home, err := dds_utils.GetHomeDirectory()
	if err != nil {
		return err
	}
	pname := pfile.GetProcessName()
	if pname == nil {
		return errors.Errorf("get process name failed the WriteProcessInfo")
	}

	file := filepath.Join(*home, "pcs", *pname)
	ok := dds_utils.IsFileExist(file)
	if ok {
		os.Remove(file)
	}
	return nil
}
