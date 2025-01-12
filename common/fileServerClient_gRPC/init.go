package fileServerClient_gRPC

import (
	"XcxcPan/Server/XcXcPanFileServer/XcXcPanFileServer"
	"XcxcPan/Server/XcXcPanFileServer/xcxcpanfileservice"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
)

var ClientOfFileServer1 = FileServerClient_1()
var ClientOfFileServer2 = FileServerClient_2()

func FileServerClient_1() XcXcPanFileServer.XcXcPanFileServiceClient {
	conn := zrpc.MustNewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: []string{"127.0.0.1:2379"},
			Key:   "fileserver1.rpc",
		},
	})
	return xcxcpanfileservice.NewXcXcPanFileService(conn)
}

func FileServerClient_2() XcXcPanFileServer.XcXcPanFileServiceClient {
	conn := zrpc.MustNewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: []string{"127.0.0.1:2379"},
			Key:   "fileserver2.rpc",
		},
	})
	return xcxcpanfileservice.NewXcXcPanFileService(conn)
}
