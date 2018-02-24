package services

import (
	"github.com/kataras/iris"
)

// 上传文件服务
type UploadService interface {
	UploadImage(ctx iris.Context) (int, error)
}

func NewUploadService() UploadService {
	return &uploadService{
	}
}

type uploadService struct {
}

// 注册
func (s *uploadService) UploadImage(ctx iris.Context) (int, error) {
	return iris.StatusOK, nil
}
