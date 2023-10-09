package oracle

import (
	"fmt"
	"github.com/892294101/dds/dbs/ddslog"
	"github.com/892294101/dds/dbs/metadata"
	"github.com/892294101/dds/dbs/spfile"
	"github.com/pingcap/errors"
	"github.com/sirupsen/logrus"
	"strings"
)

type Capture struct {
	processName, dataBaseType, processType string            //基本信息
	pfile                                  *spfile.Spfile    // 参数文件
	log                                    *logrus.Logger    // 日志记录器
	md                                     metadata.MetaData // 元数据文件
	capt                                   *CaptureTasks     // 数据捕获任务
}

func (c *Capture) readPFile() error {
	pfile, err := spfile.LoadSpfile(fmt.Sprintf("%s.desc", c.processName), spfile.UTF8, c.log, c.dataBaseType, c.processType)
	if err != nil {
		return err
	}

	if err := pfile.Production(); err != nil {
		return err
	}

	if !strings.EqualFold(*pfile.GetProcessName(), c.processName) {
		return errors.Errorf("Process name mismatch: %s", *pfile.GetProcessName())
	}
	c.pfile = pfile
	return err
}

func (c *Capture) InitConfig(processName string) error {
	c.processName = strings.ToUpper(processName)
	c.dataBaseType, c.processType = spfile.GetOracleName(), spfile.GetExtractName()

	// 初始化日志系统
	log, err := ddslog.InitDDSlog(processName)
	if err != nil {
		return err
	}
	c.log = log

	// 打开参数文件，并检查
	if err := c.readPFile(); err != nil {
		return err
	}

	// 初始化检查点元数据文件
	mds, err := metadata.InitMetaData(processName, c.dataBaseType, c.processType, c.log, metadata.LOAD)
	if err != nil {
		return err
	}
	c.md = mds

	// 初始化任务组
	tg := NewTaskGroups()
	if err := tg.InitTaskGroups(c.pfile, c.log); err != nil {
		return err
	}
	c.capt = tg

	return nil
}

func NewCapture() *Capture {
	return new(Capture)
}
