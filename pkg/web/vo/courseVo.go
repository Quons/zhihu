package vo

import "github.com/Quons/go-gin-example/models"

type CourseVo struct {
	CourseId    int64
	CourseName  string
	CourseImage string
}

func (c CourseVo) Transform(do *models.Course) *CourseVo {
	var vo CourseVo
	vo.CourseId = do.CourseID
	vo.CourseName = do.CourseName
	vo.CourseImage = do.CourseImage
	return &vo
}
