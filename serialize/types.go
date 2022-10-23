package serialize

const (
	EOF_DATA byte = 0x07
)

type Serialize interface {
	Encode(data interface{}, compress bool) ([]byte, error)
	InitBuffer() error
}
