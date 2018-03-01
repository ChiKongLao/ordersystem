package controllers

import (
	"github.com/kataras/iris"
	"github.com/chikong/ordersystem/services"
	"github.com/chikong/ordersystem/constant"
	"os"
	"io"
	"github.com/chikong/ordersystem/api/middleware/authentication"
	"github.com/chikong/ordersystem/model"
	"github.com/chikong/ordersystem/util"
	"fmt"
)

// 上传文件
type UploadController struct {
	Ctx iris.Context
	UserService services.UserService
	services.UploadService
}

func (c *UploadController) Options() (int, interface{}) {
	c.UploadImage(c.Ctx)

	return iris.StatusOK,"ok"
}



func (c *UploadController) OptionsBy() (int, interface{}) {
	userId,_ := c.Ctx.PostValueInt(constant.NameUserId)

	isOwn, err := authentication.IsOwnWithToken(c.Ctx, userId)
	if !isOwn {
		return iris.StatusUnauthorized, model.NewErrorResponse(err)
	}

	file, info, err := c.Ctx.FormFile("file")

	if err != nil {
		return iris.StatusInternalServerError, "未找到文件"
	}

	defer file.Close()

	// 创建目录
	path := fmt.Sprintf("%s/%v/",constant.PathUpload,userId)
	//fileName := fmt.Sprintf("%s-%s",strconv.FormatInt(time.Now().Unix(),10), info.Filename)
	filePathName := path + info.Filename

	if _, err := os.Stat(path); os.IsNotExist(err) {
		_ = os.MkdirAll(path, os.ModePerm)
	}
	out, err := os.OpenFile(filePathName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return iris.StatusInternalServerError, "操作文件失败"+err.Error()
	}
	defer out.Close()
	io.Copy(out, file)


	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}


	return iris.StatusOK, iris.Map{constant.Name: util.GetLocalIPWithHttp()+filePathName[1:]}
}
