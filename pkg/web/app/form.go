package app

import (
	"github.com/Quons/go-gin-example/pkg/e"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func BindAndValid(c *gin.Context, form interface{}) (int, int) {
	err := c.Bind(form)
	if err != nil {
		logrus.Errorf(err.Error())
		return http.StatusOK, e.ERROR_INVALID_PARAMS
	}

	valid := validation.Validation{}
	check, err := valid.Valid(form)
	if err != nil {
		logrus.Errorf(err.Error())
		return http.StatusOK, e.ERROR
	}
	if !check {
		MarkErrors(valid.Errors)
		return http.StatusOK, e.ERROR_INVALID_PARAMS
	}

	return http.StatusOK, e.SUCCESS
}
