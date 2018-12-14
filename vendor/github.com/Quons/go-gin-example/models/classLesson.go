package models

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"time"
)

type Classlesson struct {
	ID int64 `gorm:"primary_key;column:id"`
	//关联查询，以ForeignKey作为外键 从Course表中以AssociationForeignKey 为references进行查询
	Course    Course    `gorm:"ForeignKey:courseId;AssociationForeignKey:courseId"`
	CourseID  int64     `gorm:"column:courseId"`
	Lesson    Lesson    `gorm:"ForeignKey:lessonId;AssociationForeignKey:lessonId"`
	LessonID  int64     `gorm:"column:lessonId"`
	PhaseID   int64     `gorm:"column:phaseId"`
	StartTime time.Time `gorm:"column:startTime"`
}

func ExistClasslessonByID(id int64) (bool, error) {
	var classlesson Classlesson
	err := readDB().Select("id").Where("id = ?", id).First(&classlesson).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if classlesson.ID > 0 {
		return true, nil
	}

	return false, nil
}

func GetClasslessonTotal(maps interface{}) (int, error) {
	var count int
	if err := readDB().Model(&Classlesson{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func GetClasslessons(pageNum int, pageSize int, maps interface{}) ([]Classlesson, error) {
	var classLessons []Classlesson
	err := readDB().Order("id desc").Offset(pageNum).Limit(pageSize).Where(maps).Find(&classLessons).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return classLessons, err
	}
	return classLessons, nil
}

//获取phaseId的班级信息
func GetClasslessonsByPhase(phaseID int64) ([]Classlesson, error) {
	var classLessons []Classlesson
	err := readDB().Preload("Lesson").Preload("Lesson.Course").Preload("Course").Where("phaseId = ? and startTime < ?", phaseID, time.Now()).Find(&classLessons).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return classLessons, err
	}
	return classLessons, nil
}

func GetClasslesson(id int64) (*Classlesson, error) {
	classLesson := Classlesson{ID: id}
	err := readDB().First(&classLesson).Related(&classLesson.Course).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logrus.Errorf("%+v", err)
		return nil, err
	}
	return &classLesson, nil
}

func AddOrUpdateClasslesson(classLesson *Classlesson) error {
	if err := WriteDB().Save(classLesson).Error; err != nil {
		return err
	}
	return nil
}

func AddClasslesson(classLesson *Classlesson) error {
	return AddClasslessonTrans(WriteDB(), classLesson)
}

func AddClasslessonTrans(tx *gorm.DB, classLesson *Classlesson) error {
	if err := tx.Create(classLesson).Error; err != nil {
		return err
	}
	return nil
}

func DeleteClasslesson(id int64) error {
	if err := WriteDB().Where("id = ?", id).Delete(Classlesson{}).Error; err != nil {
		return err
	}
	return nil
}
