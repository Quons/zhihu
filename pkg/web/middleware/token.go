package middleware

import (
	"zhihu/pkg/web/model"
	"zhihu/pkg/e"
	"zhihu/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func CheckToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("studentId", "123456")
		var data interface{}

		token := c.PostForm("token")
		logrus.Info("token:", token)
		if token == "" {
			logrus.Info("empty token")
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": e.ERROR_INVALID_PARAMS,
				"msg":  e.GetMsg(e.ERROR_INVALID_PARAMS),
				"data": data,
			})
			c.Abort()
			return
		} else {

			//从用户中心拉取用户信息，并设置到
			apiStudent, err := util.GetStudentFromUserCenter(token)

			if err != nil {
				logrus.WithField("token", token).Error(err)
				c.JSON(http.StatusUnauthorized, gin.H{
					"code": e.ERROR_TOKEN_EXPIRE,
					"msg":  e.GetMsg(e.ERROR_TOKEN_EXPIRE),
					"data": data,
				})
				c.Abort()
				return
			}
			studentInfo, err := model.GetStudent(apiStudent.Data.StudentId)
			if err != nil || studentInfo.StudentID == 0 {
				logrus.WithField("token", token).Error(err)
				c.JSON(http.StatusUnauthorized, gin.H{
					"code": e.ERROR_TOKEN_EXPIRE,
					"msg":  e.GetMsg(e.ERROR_TOKEN_EXPIRE),
					"data": data,
				})
				c.Abort()
				return
			}
			c.Set(gin.AuthUserKey, studentInfo)
		}
		c.Next()
	}
}
