package spfile

const (
	IpV4Reg  = "^((0|[1-9]\\d?|1\\d\\d|2[0-4]\\d|25[0-5])\\.){3}(0|[1-9]\\d?|1\\d\\d|2[0-4]\\d|25[0-5])$"
	IpV4Port = "^([0-9]{1,4}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$"
	/*	SourceDBRegular = "(^)" +
		"(?i:(" + SourceDBType + "))  (\\s+) (((0|[1-9]\\d?|1\\d\\d|2[0-4]\\d|25[0-5])\\.){3}(0|[1-9]\\d?|1\\d\\d|2[0-4]\\d|25[0-5])) (\\s+)" +
		"(?i:(" + Port + ")) (\\s+) (\\d+) (\\s+)" +
		"(?i:(" + DataBase + ")) (*) (\\s+) " +
		"(?i:(" + Types + ")) (" + MySQL + "|" + "MariaDB" + ") (\\s+) " +
		"(?i:(" + UserId + ")) (*) (\\s+) " +
		"(?i:(" + PassWord + ")) (*) (\\s+) " +
		"($)"*/
)

const (
	MySQL    = "MySQL"
	MariaDB  = "MariaDB"
	Oracle   = "Oracle"
	Extract  = "Extract"
	Replicat = "Replicat "
)

func GetMySQLName() string {
	return MySQL
}

func GetMariaDBName() string {
	return MariaDB
}
func GetOracleName() string {
	return Oracle
}

func GetExtractName() string {
	return Extract
}

type Module interface {
	Init()
}

type Parameter interface {
	Put()
	Init()
	IsType(raw *string, dbType *string, processType *string) error
	Parse(raw *string) error
}

type Parameters interface {
	Module
	Registry() map[string]Parameter
}
