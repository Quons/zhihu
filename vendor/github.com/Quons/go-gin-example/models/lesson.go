package models

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"time"
)

type Lesson struct {
	LessonID int64 `gorm:"primary_key;column:lessonId"`
	//关联查询，以ForeignKey作为外键 从Course表中以AssociationForeignKey 为references进行查询
	//一对一，指定的是本方的外键以及对于另一方的reference
	Course     Course    `gorm:"ForeignKey:courseId;AssociationForeignKey:courseId"`
	CourseID   int64     `gorm:"column:courseId"`
	Duration   int       `gorm:"column:duration"`
	StartTime  time.Time `gorm:"column:startTime"`
	AddTime    time.Time `gorm:"column:addTime"`
	LessonNum  int       `gorm:"column:lessonNum"`
	Brief      string    `gorm:"column:brief"`
	LessonName string    `gorm:"column:lessonName"`
}

func ExistLessonByID(id int64) (bool, error) {
	var lesson Lesson
	err := readDB().Select("lessonId").Where("lessonId = ?", id).First(&lesson).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if lesson.LessonID > 0 {
		return true, nil
	}

	return false, nil
}

func GetLessonTotal(maps interface{}) (int, error) {
	var count int
	if err := readDB().Model(&Lesson{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func GetLessons(pageNum int, pageSize int, maps interface{}) ([]Lesson, error) {
	var lessons []Lesson
	err := readDB().Order("lessonId desc").Offset(pageNum).Limit(pageSize).Where(maps).Find(&lessons).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return lessons, err
	}
	return lessons, nil
}

func GetLesson(id int64) (*Lesson, error) {
	lesson := Lesson{LessonID: id}
	err := readDB().First(&lesson).Related(&lesson.Course).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logrus.Errorf("%+v", err)
		return nil, err
	}
	return &lesson, nil
}

func AddOrUpdateLesson(lesson *Lesson) error {
	if err := WriteDB().Save(lesson).Error; err != nil {
		return err
	}
	return nil
}

func AddLesson(lesson *Lesson) error {
	return AddLessonTrans(WriteDB(), lesson)
}

func AddLessonTrans(tx *gorm.DB, lesson *Lesson) error {
	if err := tx.Create(lesson).Error; err != nil {
		return err
	}
	return nil
}

func DeleteLesson(id int64) error {
	if err := WriteDB().Where("lessonId = ?", id).Delete(Lesson{}).Error; err != nil {
		return err
	}
	return nil
}
