package api

import (
	"github.com/gin-gonic/gin"

	"zhihu/pkg/web/app"
	"zhihu/pkg/e"
	"zhihu/pkg/web/upload"
	"github.com/sirupsen/logrus"
)

// @Summary 上传图片
// @Produce  json
// @Param image post file true "图片文件"
// @Success 200 {string} json "{"code":200,"data":{"image_save_url":"upload/images/96a.jpg", "image_url": "http://..."}"
// @Router /api/v1/tags/import [post]
func UploadImage(c *gin.Context) {
	appG := app.Gin{C: c}
	file, image, err := c.Request.FormFile("image")
	if err != nil {
		logrus.Warn(err)
		appG.Response(nil, e.ERROR)
		return
	}

	if image == nil {
		appG.Response(nil, e.ERROR_INVALID_PARAMS)
		return
	}

	imageName := upload.GetImageName(image.Filename)
	fullPath := upload.GetImageFullPath()
	savePath := upload.GetImagePath()
	src := fullPath + imageName

	if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(file) {
		appG.Response(nil, e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT)
		return
	}

	err = upload.CheckImage(fullPath)
	if err != nil {
		logrus.Warn(err)
		appG.Response(nil, e.ERROR_UPLOAD_CHECK_IMAGE_FAIL)
		return
	}

	if err := c.SaveUploadedFile(image, src); err != nil {
		logrus.Warn(err)
		appG.Response(nil, e.ERROR_UPLOAD_SAVE_IMAGE_FAIL)
		return
	}

	appG.Response(map[string]string{
		"image_url":      upload.GetImageFullUrl(imageName),
		"image_save_url": savePath + imageName,
	}, e.SUCCESS)
}
