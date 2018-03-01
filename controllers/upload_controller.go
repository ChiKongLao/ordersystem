package controllers

import (
	"github.com/kataras/iris"
	"github.com/chikong/ordersystem/services"
	"github.com/chikong/ordersystem/constant"
	"github.com/chikong/ordersystem/api/middleware/authentication"
	"github.com/chikong/ordersystem/model"
	"fmt"
)

// 上传文件
type UploadController struct {
	Ctx         iris.Context
	UserService services.UserService
	services.UploadService
}

func (c *UploadController) Options() (int, interface{}) {
	userId, _ := c.Ctx.PostValueInt(constant.NameUserId)
	if userId == 0 {
		return iris.StatusBadRequest, model.NewErrorResponseWithMsg("用户不能为空")

	}

	isOwn, err := authentication.IsOwnWithToken(c.Ctx, userId)
	if !isOwn {
		return iris.StatusUnauthorized, model.NewErrorResponse(err)
	}

	file, info, err := c.Ctx.FormFile("file")
	if err != nil {
		return iris.StatusBadRequest, model.NewErrorResponseWithMsg("文件不能为空")
	}

	defer file.Close()

	status, err := c.UploadImage(file, info.Filename, userId)

	if err != nil {
		return status, model.NewErrorResponse(err)
	}

	return iris.StatusOK, iris.Map{
		constant.Name: fmt.Sprintf("%s/%s", constant.UrlQiNiuUrl, info.Filename),
	}
}
