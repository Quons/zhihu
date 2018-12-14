package v1

import (
	"zhihu/pkg/web/model"
	"zhihu/pkg/web/app"
	"zhihu/pkg/e"
	"zhihu/pkg/web/handler/course_handler"
	"zhihu/pkg/web/vo"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// @Tags 课程
// @Summary 获取单个课程
// @Description 获取单个课程description
// @Produce  json
// @accept application/x-www-form-urlencoded
// @Param token formData string true "用户token"
// @Param courseId formData int true "课程ID"
// @Success 200 {object} vo.CourseVo
// @Failure 10000 {string} json "{"code":10000,"data":{},"msg":"服务器错误"}"
// @Failure 20000 {string} json "{"code":20000,"data":{},"msg":"参数错误"}"
// @Router /api/v1/getCourse [post]
func GetCourse(c *gin.Context) {
	appG := app.Gin{C: c}
	course := course_handler.Course{}
	err := c.ShouldBind(&course)
	if err != nil {
		log.Info(err)
		appG.Response(nil, e.ERROR_INVALID_PARAMS)
		return
	}
	log.WithField("courseId", course.CourseID).Info()
	//获取student信心
	studentInfo := c.MustGet(gin.AuthUserKey).(*model.Student)
	log.WithField("studentInfo", studentInfo).Info()
	courseDO, err := course.Get()
	if err != nil {
		appG.Response(nil, e.ERROR_SERVER_ERROR)
		return
	}
	if courseDO.CourseID == 0 {
		appG.Response(nil, e.ERROR_DATA_ERROR)
		return
	}
	appG.Response(vo.CourseVo{}.Transform(courseDO), e.SUCCESS)
}
