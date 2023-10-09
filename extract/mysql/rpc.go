package oramysql

import (
	"context"
	"github.com/892294101/ddsrpc/pcb"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"strconv"
)

const (
	STOP = iota + 1
	FORCE
	STATISTICS
	DETAILED
)

type RpcProcess struct {
	pro    *ExtractEvent
	server *grpc.Server
	log    *logrus.Logger
}

func (r *RpcProcess) Stop(ctx context.Context, command *pcb.StopCommand) (*pcb.StopReply, error) {
	r.log.Infof("%v", "Stop signal processing")
	r.pro.stopSrv.power <- STOP
	return &pcb.StopReply{}, nil
}

func (r *RpcProcess) CloseRpc() {
	r.log.Infof("turn off rpc server")
	r.server.Stop()
}

func NewRpc() *RpcProcess {
	return new(RpcProcess)
}

func (r *RpcProcess) StartRpcServer(p *ExtractEvent, log *logrus.Logger) (err error) {
	if p == nil {
		log.Fatalf("rpc must pass in entity")
	}
	r.pro = p   // 被控制的进程
	r.log = log // 日志记录器

	lis, err := net.Listen("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(r.pro.rpcPort)))
	if err != nil {
		return errors.Errorf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	s := grpc.NewServer(opts...)

	// 注册Greeter服务
	pcb.RegisterGreeterServer(s, r)
	// 往grpc服务端注册反射服务
	reflection.Register(s)
	r.server = s // rpc server

	// 启动grpc服务
	r.pro.isReady = true
	if err := s.Serve(lis); err != nil {
		r.pro.isReady = false
		return errors.Errorf("failed to serve: %v", err)
	}

	return nil
}
