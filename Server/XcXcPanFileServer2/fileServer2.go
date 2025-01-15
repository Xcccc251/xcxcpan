package main

import (
	"XcxcPan/Server/XcXcPanFileServer2/XcXcPanFileServer"
	"XcxcPan/Server/XcXcPanFileServer2/internal/config"
	"XcxcPan/Server/XcXcPanFileServer2/internal/server"
	"XcxcPan/Server/XcXcPanFileServer2/internal/svc"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f2", "etc/fileserver.yaml", "the config file")

func main() {
	flag.Parse()
	var c2 config.Config
	conf.MustLoad(*configFile, &c2)
	ctx2 := svc.NewServiceContext(c2)

	s2 := zrpc.MustNewServer(c2.RpcServerConf, func(grpcServer *grpc.Server) {
		XcXcPanFileServer.RegisterXcXcPanFileServiceServer(grpcServer, server.NewXcXcPanFileServiceServer(ctx2))

		if c2.Mode == service.DevMode || c2.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s2.Stop()

	fmt.Printf("Starting rpc server2 at %s...\n", c2.ListenOn)
	s2.Start()

}
