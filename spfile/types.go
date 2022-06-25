package spfile

const (
	ProcessType     = "PROCESS"
	ProcessRegular  = "(^)(?i:(" + ProcessType + "))(\\s+)((?:[A-Za-z0-9_]){4,12})($)"
	SourceDBType    = "SOURCEDB"
	SourceDBRegular = "(^)(?i:(" + SourceDBType + "))(\\s+)((?:[A-Za-z0-9_]){4,12})($)"
)

type Module interface {
	Init()
}

type Parameter interface {
	Put()
	Parse(raw string) error
}

type Parameters interface {
	Module
	Registry() map[string]Parameter
}
