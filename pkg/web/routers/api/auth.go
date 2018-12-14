package api

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"

	"zhihu/pkg/web/app"
	"zhihu/pkg/e"
	"zhihu/pkg/util"
	"zhihu/pkg/web/handler/auth_handler"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

/*GetAuth 验证信息*/
func GetAuth(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}

	username := c.Query("username")
	password := c.Query("password")

	a := auth{Username: username, Password: password}
	ok, _ := valid.Valid(&a)

	if !ok {
		app.MarkErrors(valid.Errors)
		appG.Response(nil, e.ERROR_INVALID_PARAMS)
		return
	}

	authService := auth_handler.Auth{Username: username, Password: password}
	isExist, err := authService.Check()
	if err != nil {
		appG.Response(nil, e.ERROR_AUTH_CHECK_TOKEN_FAIL)
		return
	}

	if !isExist {
		appG.Response(nil, e.ERROR_AUTH)
		return
	}

	token, err := util.GenerateToken(username, password)
	if err != nil {
		appG.Response(nil, e.ERROR_AUTH_TOKEN)
		return
	}

	appG.Response(map[string]string{
		"token": token,
	}, e.SUCCESS)
}
