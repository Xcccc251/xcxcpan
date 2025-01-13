package logic

import (
	Server_MinIO "XcxcPan/Server/XcXcPanFileServer/minIO"
	"context"

	"XcxcPan/Server/XcXcPanFileServer/XcXcPanFileServer"
	"XcxcPan/Server/XcXcPanFileServer/internal/svc"

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
	err := Server_MinIO.DelChunk(in.FileName, int(in.Server))
	if err != nil {
		return nil, err
	}
	return &XcXcPanFileServer.DelChunkResponse{
		Message: "删除成功",
	}, nil
}
