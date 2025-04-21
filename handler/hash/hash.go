package hash

import (
	"github.com/gin-gonic/gin"

	"learning-assistant/handler/aerrors"
	"learning-assistant/handler/basic"
	"learning-assistant/service/hash"
)

// BindImageHash
// @Summary 绑定图片hash信息
// @Tag Hash(图片hash)
// @Param path query string true "仓库名称"
// @Param file formData file true "上传图片文件"
// @Success 200 {object} basic.Resp{}
// @Router /image/hash/bind [POST]
func BindImageHash(c *gin.Context) {
	// 从请求中获取文件
	_, header, err := c.Request.FormFile("file")
	repo := c.Query("path")
	if err != nil {
		basic.RequestParamsFailure(c)
		return
	}
	err = hash.BindImageHash(c, header, repo)
	if err != nil {
		basic.RequestFailureWithCode(c, err.Error(), aerrors.NormalError)
		return
	}
	basic.Success(c, nil)
	return
}

// SimilarImage
// @Summary 查询图片
// @Tag Hash(图片hash)
// @Param path query string true "仓库名称"
// @Param file formData file true "上传图片文件"
// @Success 200 {object} basic.Resp{}
// @Router /image/hash/similar [GET]
func SimilarImage(c *gin.Context) {
	_, header, err := c.Request.FormFile("file")
	repo := c.Query("path")
	if err != nil {
		basic.RequestParamsFailure(c)
		return
	}
	image, err := hash.QuerySimilarImage(c, header, repo)
	if err != nil {
		basic.RequestFailureWithCode(c, err.Error(), aerrors.NormalError)
		return
	}
	basic.Success(c, image)
	return
}
