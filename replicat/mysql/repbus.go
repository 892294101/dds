package mysql

import (
	"github.com/892294101/dds/replicat/conn"
	"github.com/892294101/dds/spfile"
)

type replicat struct {
	pfile spfile.Spfile
	conn  conn.Connector
}

func (r *replicat) InitEnv() {

}

func (r *replicat) DetectDataType() {

}
