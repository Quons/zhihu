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

func TestGetLesson(t *testing.T) {
	a, err := GetLesson(1)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%+v", a)
	assert.Equal(t, int64(1), a.LessonID)
}

func TestExistLessonByID(t *testing.T) {
	c, err := ExistLessonByID(1)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%+v", c)
	assert.Equal(t, true, c)
}

func TestGetLessonTotal(t *testing.T) {
	c, err := GetLessonTotal(map[string]interface{}{})
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%+v", c)
	assert.Equal(t, 24, c)
}

func TestGetLessons(t *testing.T) {
	Lessons, err := GetLessons(0, 11, map[string]interface{}{})
	if err != nil {
		t.Error(err)
		return
	}
	for _, Lesson := range Lessons {
		t.Logf("%+v", Lesson)
	}
	assert.Equal(t, 11, len(Lessons))
}

func TestAddLesson(t *testing.T) {
	Lesson := &Lesson{LessonName: "testLesson"}
	err := AddLesson(Lesson)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestEditLesson(t *testing.T) {
	Lesson := &Lesson{LessonID: 20, LessonName: "testLessonsss"}
	err := AddOrUpdateLesson(Lesson)
	if err != nil {
		t.Error(err)
		return
	}
}
