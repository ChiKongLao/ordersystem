package services

import (
	"github.com/kataras/iris"
	"github.com/qiniu/api.v7/storage"
	"github.com/qiniu/api.v7/auth/qbox"
	"fmt"
	"context"
	"github.com/chikong/ordersystem/constant"
	"sync"
	"mime/multipart"
	"os"
	"io"
	"github.com/sirupsen/logrus"
	"github.com/kataras/iris/core/errors"
)

// 上传文件服务
type UploadService interface {
	UploadImage(file multipart.File, fileName string, userId int) (int, error)
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
	mac := qbox.NewMac(constant.KeyQiNiuAccessKey,constant.KeyQiNiuSecretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	cfg.Zone = &storage.ZoneHuanan
	cfg.UseHTTPS = false
	cfg.UseCdnDomains = false

	return storage.NewFormUploader(&cfg),upToken
}

// 注册
func (s *uploadService) UploadImage(file multipart.File, fileName string, userId int) (int, error) {
	mOnce.Do(func() {
		initService()
	})

	// 创建目录
	path := fmt.Sprintf("%s/%v/",constant.PathUpload,userId)
	filePath := path + fileName

	if _, err := os.Stat(path); os.IsNotExist(err) {
		_ = os.MkdirAll(path, os.ModePerm)
	}
	out, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logrus.Errorf("读取文件失败. "+err.Error())
		return iris.StatusInternalServerError, errors.New("读取文件失败")
	}
	defer out.Close()
	io.Copy(out, file)

	if err != nil {
		logrus.Errorf("操作文件失败. "+err.Error())
		return iris.StatusInternalServerError, errors.New("操作文件失败")
	}

	uploader, upToken := initService()
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
	}

	err = uploader.PutFile(context.Background(), &ret, upToken, fileName, filePath, &putExtra)
	if err != nil {
		logrus.Errorf("上传文件失败. "+err.Error())
		return iris.StatusInternalServerError, errors.New("上传文件失败")
	}
	os.Remove(filePath) // 上传成功, 删除本地缓存的图片
	return iris.StatusOK, nil
}
