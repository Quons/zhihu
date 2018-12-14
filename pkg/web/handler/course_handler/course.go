package course_handler

import (
	"github.com/Quons/go-gin-example/models"
	log "github.com/sirupsen/logrus"
)

type Course struct {
	CourseID int64 `form:"courseId" json:"courseId" binding:"required,gte=0,lte=4"`
}

func (c *Course) Get() (*models.Course, error) {
	course, err := models.GetCourse(c.CourseID)
	if err != nil {
		log.WithField("courseId", c.CourseID).Error(err)
		return course, err
	}
	return course, nil
}

