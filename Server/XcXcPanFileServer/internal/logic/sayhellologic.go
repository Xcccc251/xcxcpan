package logic

import (
	"XcxcPan/Server/XcXcPanFileServer/XcXcPanFileServer"
	"XcxcPan/Server/XcXcPanFileServer/internal/svc"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type SayHelloLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSayHelloLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SayHelloLogic {
	return &SayHelloLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SayHelloLogic) SayHello(in *XcXcPanFileServer.SayHelloRequest) (*XcXcPanFileServer.SayHelloResponse, error) {

	return &XcXcPanFileServer.SayHelloResponse{
		Message: "Hello! " + in.Name,
	}, nil
}
