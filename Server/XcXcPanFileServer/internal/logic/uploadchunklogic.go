package logic

import (
	"XcxcPan/Server/common/helper"
	Server_MinIO "XcxcPan/Server/common/minIO"
	"XcxcPan/common/define"

	"context"

	"XcxcPan/Server/XcXcPanFileServer/XcXcPanFileServer"
	"XcxcPan/Server/XcXcPanFileServer/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadChunkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUploadChunkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadChunkLogic {
	return &UploadChunkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UploadChunkLogic) UploadChunk(in *XcXcPanFileServer.UploadChunkRequest) (*XcXcPanFileServer.UploadChunkResponse, error) {
	if Server_MinIO.CheckChunkExists(in.FileName, define.Server1) {
		return &XcXcPanFileServer.UploadChunkResponse{
			Message: "文件已存在",
		}, nil
	}

	file, err := Server_Helper.BytesToFile(in.Data)
	if err != nil {
		return nil, err
	}

	err = Server_MinIO.UploadChunk(in.FileName, file, define.Server1)
	if err != nil {
		return nil, err
	}

	return &XcXcPanFileServer.UploadChunkResponse{
		Message: "上传成功",
	}, nil
}
