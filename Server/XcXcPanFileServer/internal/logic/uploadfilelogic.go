package logic

import (
	Server_Helper "XcxcPan/Server/XcXcPanFileServer/helper"
	Server_MinIO "XcxcPan/Server/XcXcPanFileServer/minIO"
	"context"

	"XcxcPan/Server/XcXcPanFileServer/XcXcPanFileServer"
	"XcxcPan/Server/XcXcPanFileServer/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUploadFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadFileLogic {
	return &UploadFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UploadFileLogic) UploadFile(in *XcXcPanFileServer.UploadFileRequest) (*XcXcPanFileServer.UploadFileResponse, error) {
	file, err := Server_Helper.BytesToFile(in.Data)
	if err != nil {
		return nil, err
	}
	_, err = Server_MinIO.UploadImage(in.FileName, file, 1)
	if err != nil {
		return nil, err
	}

	return &XcXcPanFileServer.UploadFileResponse{
		Message: "上传成功",
	}, nil
}
