package logic

import (
	Server_Helper "XcxcPan/Server/XcXcPanFileServer/helper"
	Server_MinIO "XcxcPan/Server/XcXcPanFileServer/minIO"
	"context"

	"XcxcPan/Server/XcXcPanFileServer/XcXcPanFileServer"
	"XcxcPan/Server/XcXcPanFileServer/internal/svc"

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
	file, err := Server_MinIO.DownloadChunk(in.FileName, int(in.Server))
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
