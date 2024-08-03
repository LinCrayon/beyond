// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	"github.com/LinCrayon/beyond/application/article/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/publish",
				Handler: PublishHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/upload/cover",
				Handler: UploadCoverHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		rest.WithPrefix("/v1/article"),
	)
}
