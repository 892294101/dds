package oracle

import "database/sql"

// 用于接收Oracle 11.x 版本V$LOGMNR_CONTENTS视图的内容
type logContentV11 struct {
	Thread        int           `json:"ThreadID"`  //NUMBER
	Sequence      sql.NullInt64 `json:"Sequence"`  //NUMBER
	Scn           sql.NullInt64 `json:"Scn"`       //NUMBER
	StartScn      sql.NullInt64 `json:"StartScn"`  //NUMBER
	CommitScn     sql.NullInt64 `json:"CommitScn"` //NUMBER
	Xid           string        `json:"Xid"`
	OperationCode sql.NullInt64 `json:"OperationCode"`
	Status        string        `json:"Status"`
	SegOwner      string        `json:"SegOwner"`
	TableName     string        `json:"TableName"`
	UserName      string        `json:"UserName"`
	SegName       string        `json:"SegName"`
	SegType       string        `json:"SegType"`
	SqlRedo       string        `json:"SqlRedo"`
	RowId         string        `json:"RowId"`
	TableSpace    string        `json:"TableSpace"`
	Rollback      int           `json:"Rollback"`
	RsId          string        `json:"RsId"`
	Ssn           sql.NullInt64 `json:"Ssn"`
	Csf           sql.NullInt64 `json:"Csf"`
	RbaSqn        sql.NullInt64 `json:"RbaSqn"`
	RbaBlk        sql.NullInt64 `json:"RbaBlk"`
}
