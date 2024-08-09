package logic

import (
	"context"
	"fmt"
	"github.com/LinCrayon/beyond/application/article/api/internal/code"
	"net/http"
	"time"

	"github.com/LinCrayon/beyond/application/article/api/internal/svc"
	"github.com/LinCrayon/beyond/application/article/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

const maxFileSize = 10 << 20 // 10MB

type UploadCoverLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadCoverLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadCoverLogic {
	return &UploadCoverLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UploadCover 封面上传
func (l *UploadCoverLogic) UploadCover(req *http.Request) (resp *types.UploadCoverResponse, err error) {
	//解析了 HTTP 请求中的 multipart 表单，用于处理文件上传
	_ = req.ParseMultipartForm(maxFileSize)
	//接收form表单传过来的文件，文件名cover
	file, handler, err := req.FormFile("cover")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	//获取存储桶（Bucket）
	bucket, err := l.svcCtx.OssClient.Bucket(l.svcCtx.Config.Oss.BucketName)
	if err != nil {
		logx.Errorf("get bucket failed, err: %v", err)
		return nil, code.GetBucketErr
	}
	//生成文件名
	objectKey := genFilename(handler.Filename)
	//发送文件到oss
	err = bucket.PutObject(objectKey, file)
	if err != nil {
		logx.Errorf("put object failed, err: %v", err)
		return nil, code.PutBucketErr
	}
	return &types.UploadCoverResponse{CoverUrl: genFileURL(objectKey)}, nil
}

// 生成文件名
func genFilename(filename string) string {
	return fmt.Sprintf("%d_%s", time.Now().UnixMilli(), filename)
}

func genFileURL(objectKey string) string {
	return fmt.Sprintf("https://beyond-article01.oss-cn-hangzhou.aliyuncs.com/%s", objectKey)
}
