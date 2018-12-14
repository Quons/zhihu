package app

import (
	"github.com/gin-gonic/gin"

	"github.com/Quons/go-gin-example/pkg/e"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Gin struct {
	C *gin.Context
}

/*返回实体，code[0]为业务返回码，code[1]为http返回码*/
func (g *Gin) Response(data interface{}, code ...int) {
	if len(code) < 1 {
		logrus.Error("invalid code param")
		g.C.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR_SERVER_ERROR,
			"msg":  e.GetMsg(e.ERROR_SERVER_ERROR),
			"data": data,
		})
	}

	if len(code) == 1 {
		g.C.JSON(http.StatusOK, gin.H{
			"code": code[0],
			"msg":  e.GetMsg(code[0]),
			"data": data,
		})
	}

	if len(code) == 2 {
		g.C.JSON(code[1], gin.H{
			"code": code[0],
			"msg":  e.GetMsg(code[0]),
			"data": data,
		})
	}
	return
}
