package services

import (
	"github.com/kataras/iris"
	"github.com/qiniu/api.v7/storage"
	"github.com/qiniu/api.v7/auth/qbox"
	"fmt"
	"context"
	"github.com/chikong/ordersystem/constant"
	"sync"
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

var mOnce sync.Once

// 初始化上传服务
func initService() (*storage.FormUploader, string){
	putPolicy := storage.PutPolicy{
		Scope: constant.KeyQiNiuBucket,
	}
	mac := qbox.NewMac(constant.KeyQiNiuAccessKey,constant.KeyQiNiuAccessKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	cfg.Zone = &storage.ZoneHuanan
	cfg.UseHTTPS = false
	cfg.UseCdnDomains = false

	return storage.NewFormUploader(&cfg),upToken
}

// 注册
func (s *uploadService) UploadImage(ctx iris.Context) (int, error) {
	mOnce.Do(func() {
		initService()
	})
	localFile := "./1.jpg"
	key := "1.jpg"

	uploader, upToken := initService()
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
	}

	err := uploader.PutFile(context.Background(), &ret, upToken, key, localFile, &putExtra)
	if err != nil {
		fmt.Println(err)
		return iris.StatusInternalServerError,nil
	}
	fmt.Println(ret.Key, ret.Hash)
	return iris.StatusOK, nil
}
