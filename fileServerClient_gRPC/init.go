package fileServerClient_gRPC

import (
	"XcxcPan/Server/XcXcPanFileServer/XcXcPanFileServer"
	"XcxcPan/Server/XcXcPanFileServer/xcxcpanfileservice"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
)

var ClientOfFileServer1 = FileServerClient_1()
var ClientOfFileServer2 = FileServerClient_2()

var ClientMap = map[int]XcXcPanFileServer.XcXcPanFileServiceClient{
	1: ClientOfFileServer1,
	2: ClientOfFileServer2,
}

func GetClientById(id int) XcXcPanFileServer.XcXcPanFileServiceClient {
	return ClientMap[id]
}
func FileServerClient_1() XcXcPanFileServer.XcXcPanFileServiceClient {
	conn := zrpc.MustNewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: []string{"1.94.166.62:2379"},
			Key:   "fileserver1.rpc",
		},
	})
	return xcxcpanfileservice.NewXcXcPanFileService(conn)
}

func FileServerClient_2() XcXcPanFileServer.XcXcPanFileServiceClient {
	conn := zrpc.MustNewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: []string{"1.94.166.62:2379"},
			Key:   "fileserver2.rpc",
		},
	})
	return xcxcpanfileservice.NewXcXcPanFileService(conn)
}
