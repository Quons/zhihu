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

func TestGetClassLesson(t *testing.T) {
	a, err := GetLesson(1)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%+v", a)
	assert.Equal(t, int64(1), a.LessonID)
}

func TestGetClasslessonsByPhase(t *testing.T) {
	classLessons, err := GetClasslessonsByPhase(2)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%+v", classLessons)
	for _, value := range classLessons {
		t.Logf("%+v", value.Lesson)
		t.Logf("startTime:%v", value.StartTime)
	}
}

func TestExistClassLessonByID(t *testing.T) {
	c, err := ExistLessonByID(1)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%+v", c)
	assert.Equal(t, true, c)
}

func TestGetClassLessonTotal(t *testing.T) {
	c, err := GetClasslessonTotal(map[string]interface{}{})
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%+v", c)
	assert.Equal(t, 79, 79)
}

func TestGetClassLessons(t *testing.T) {
	Lessons, err := GetClasslessons(0, 11, map[string]interface{}{})
	if err != nil {
		t.Error(err)
		return
	}
	for _, Lesson := range Lessons {
		t.Logf("%+v", Lesson)
	}
	assert.Equal(t, 10, len(Lessons))
}

func TestAddClassLesson(t *testing.T) {
	Lesson := &Lesson{LessonName: "testLesson"}
	err := AddLesson(Lesson)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestEditClassLesson(t *testing.T) {
	Lesson := &Classlesson{LessonID: 20, PhaseID: 22}
	err := AddOrUpdateClasslesson(Lesson)
	if err != nil {
		t.Error(err)
		return
	}
}
