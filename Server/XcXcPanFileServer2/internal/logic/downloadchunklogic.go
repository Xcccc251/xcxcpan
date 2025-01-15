package logic

import (
	"XcxcPan/Server/XcXcPanFileServer2/XcXcPanFileServer"
	"XcxcPan/Server/XcXcPanFileServer2/internal/svc"
	"XcxcPan/Server/common/helper"
	Server_MinIO "XcxcPan/Server/common/minIO"
	"XcxcPan/common/define"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
)

type DownloadChunkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDownloadChunkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DownloadChunkLogic {
	return &DownloadChunkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DownloadChunkLogic) DownloadChunk(in *XcXcPanFileServer.DownloadChunkRequest) (*XcXcPanFileServer.DownloadChunkResponse, error) {
	file, err := Server_MinIO.DownloadChunk(in.FileName, define.Server2)
	if err != nil {
		return nil, err
	}
	data, err := Server_Helper.FileToBytes(file)
	if err != nil {
		return nil, err
	}

	return &XcXcPanFileServer.DownloadChunkResponse{
		Data: data,
	}, nil
}
