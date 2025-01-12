package main

import (
	"XcxcPan/Server/XcXcPanFileServer/XcXcPanFileServer"
	"XcxcPan/Server/XcXcPanFileServer/internal/config"
	"XcxcPan/Server/XcXcPanFileServer/internal/server"
	"XcxcPan/Server/XcXcPanFileServer/internal/svc"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"sync"
)

var configFile1 = flag.String("f1", "etc/fileserver1.yaml", "the config file")
var configFile2 = flag.String("f2", "etc/fileserver2.yaml", "the config file")

func main() {
	flag.Parse()
	wg := sync.WaitGroup{}
	wg.Add(2)
	//模拟多台服务
	go func() {
		defer wg.Done()
		var c config.Config
		conf.MustLoad(*configFile1, &c)
		ctx := svc.NewServiceContext(c)

		s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
			XcXcPanFileServer.RegisterXcXcPanFileServiceServer(grpcServer, server.NewXcXcPanFileServiceServer(ctx))

			if c.Mode == service.DevMode || c.Mode == service.TestMode {
				reflection.Register(grpcServer)
			}
		})
		defer s.Stop()

		fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
		s.Start()
	}()

	go func() {
		defer wg.Done()
		var c2 config.Config
		conf.MustLoad(*configFile2, &c2)
		ctx2 := svc.NewServiceContext(c2)

		s2 := zrpc.MustNewServer(c2.RpcServerConf, func(grpcServer *grpc.Server) {
			XcXcPanFileServer.RegisterXcXcPanFileServiceServer(grpcServer, server.NewXcXcPanFileServiceServer(ctx2))

			if c2.Mode == service.DevMode || c2.Mode == service.TestMode {
				reflection.Register(grpcServer)
			}
		})
		defer s2.Stop()

		fmt.Printf("Starting rpc server at %s...\n", c2.ListenOn)
		s2.Start()
	}()

	//// 捕获系统信号
	//quit := make(chan os.Signal, 1)
	//signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	//
	//fmt.Println("Servers are running... Press Ctrl+C to exit.")
	//<-quit // 等待信号

	//fmt.Println("Shutting down servers...")
	wg.Wait() // 等待所有服务退出
	fmt.Println("Servers stopped gracefully.")

}
