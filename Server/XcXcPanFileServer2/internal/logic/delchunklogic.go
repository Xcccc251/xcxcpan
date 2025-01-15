package logic

import (
	"XcxcPan/Server/XcXcPanFileServer2/XcXcPanFileServer"
	"XcxcPan/Server/XcXcPanFileServer2/internal/svc"
	Server_MinIO "XcxcPan/Server/common/minIO"
	"XcxcPan/common/define"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
)

type DelChunkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelChunkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelChunkLogic {
	return &DelChunkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelChunkLogic) DelChunk(in *XcXcPanFileServer.DelChunkRequest) (*XcXcPanFileServer.DelChunkResponse, error) {
	err := Server_MinIO.DelChunk(in.FileName, define.Server2)
	if err != nil {
		return nil, err
	}
	return &XcXcPanFileServer.DelChunkResponse{
		Message: "删除成功",
	}, nil
}
