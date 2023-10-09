package stream

type StreamForMySQL struct {
	stream chan struct{}
}

type StreamForOracle struct {
	stream chan struct{}
}

