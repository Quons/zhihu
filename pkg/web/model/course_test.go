package model

import (
	"zhihu/pkg/setting"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	setting.Setup("dev")
	Setup()
}

func TestGetCourse(t *testing.T) {
	a, err := GetCourse(1)
	if err != nil {
		t.Error(err)
		return
	}
	for _, value := range a.Lessons {
		t.Log(value.LessonName)
	}
	for _, value := range a.LessonSections {
		t.Log(value.SectionName)

	}
	assert.Equal(t, int64(1), a.CourseID)
}

func TestExistCourseByID(t *testing.T) {
	c, err := ExistCourseByID(1)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%+v", c)
	assert.Equal(t, true, c)
}

func TestGetCourseTotal(t *testing.T) {
	c, err := GetCourseTotal(map[string]interface{}{"status": 1})
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%+v", c)
	assert.Equal(t, 10, c)
}

func TestGetCourses(t *testing.T) {
	courses, err := GetCourses(0, 11, map[string]interface{}{"status": 1})
	if err != nil {
		t.Error(err)
		return
	}
	for _, course := range courses {
		t.Logf("%+v", course)
	}
	assert.Equal(t, 10, len(courses))
}

func TestAddCourse(t *testing.T) {
	course := &Course{CourseName: "testCourse", CourseImage: "httpssss"}
	err := AddCourse(course)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(course.CourseID)
}

func TestEditCourse(t *testing.T) {
	course, err := GetCourse(20)
	if err != nil {
		t.Error(err)
		return
	}
	course.CourseID = 21
	course.CourseName = "hiahiahi"
	course.Status = 0
	err = AddOrUpdateCourse(course)
	if err != nil {
		t.Error(err)
		return
	}
}
