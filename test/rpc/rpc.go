package main

import (
	"fmt"
	rpc "github.com/892294101/dds/grpc"
)

func main() {
	c, err := rpc.NewRpcClient()

	if err != nil {
		fmt.Println("NewRpc error", err)
	}
	err = c.Stop()
	if err != nil {
		fmt.Println("Stop error", err)

	}
}
