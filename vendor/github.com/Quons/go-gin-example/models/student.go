package models

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"time"
)

type Student struct {
	StudentID             int64     `gorm:"primary_key;column:studentId"`
	Mobile                string    `gorm:"column:mobile"`
	HealthPoint           int       `gorm:"column:healthPoint"`
	MaxHealthPoint        int       `gorm:"column:maxHealthPoint"`
	LifeCount             int       `gorm:"column:lifeCount"`
	HealthPointUpdateTime time.Time `gorm:"column:healthPointUpdateTime"`
	AddTime               time.Time `gorm:"column:addTime"`
	LifeCountUpdateTime   time.Time `gorm:"column:lifeCountUpdateTime"`
	UserId                int64     `gorm:"column:userId"`
	UserName              string    `gorm:"column:userName"`
	Channel               string    `gorm:"column:channel"`
	Platform              string    `gorm:"column:platform"`
	SystemInfo            string    `gorm:"column:systemInfo"`
	Version               int       `gorm:"column:version"`
	LearningCard          int       `gorm:"column:learningCard"`
	LastOpenTime          time.Time `gorm:"column:lastOpenTime"`
	Coin                  int       `gorm:"column:coin"`
	AbilityValue          int       `gorm:"column:abilityValue"`
	AbilityLevel          int       `gorm:"column:abilityLevel"`
	HeadPhoto             string    `gorm:"column:headPhoto"`
	Sex                   int       `gorm:"column:sex"`
	IsEdit                int       `gorm:"column:isEdit"`
	ThroughCount          int       `gorm:"column:throughCount"`
	ThroughTimes          int       `gorm:"column:throughTimes"`
	TotalCodeRows         int       `gorm:"column:totalCodeRows"`
	TotalCodeTime         int       `gorm:"column:totalCodeTime"`
	GuideProgress         string    `gorm:"column:guideProgress"`
	LinkMobile            string    `gorm:"column:linkMobile"`
	Age                   int       `gorm:"column:age"`
	Source                string    `gorm:"column:source"`
}

/*// 设置User的表名为`profiles`
func (Course) TableName() string {
	return "cb_course"
}*/

func ExistStudentByID(id int64) (bool, error) {
	var student Student
	err := readDB().Select("studentId").Where("studentId = ?", id).First(&student).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if student.StudentID > 0 {
		return true, nil
	}

	return false, nil
}

func GetStudentTotal(maps interface{}) (int, error) {
	var count int
	if err := readDB().Model(&Student{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func GetStudents(pageNum int, pageSize int, maps interface{}) ([]Student, error) {
	var student []Student
	err := readDB().Order("studentId desc").Offset(pageNum).Limit(pageSize).Where(maps).Find(&student).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return student, err
	}
	return student, nil
}

func GetStudent(id int64) (*Student, error) {
	var student Student
	student.StudentID = id
	err := readDB().First(&student).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logrus.Errorf("%+v", err)
		return nil, err
	}
	return &student, nil
}

//id为0时新增，否则更新
func AddOrUpdateStudent(model *Student) error {
	if err := WriteDB().Save(model).Error; err != nil {
		return err
	}
	return nil
}

func AddStudent(model *Student) error {
	return AddStudentTrans(WriteDB(), model)
}

func AddStudentTrans(tx *gorm.DB, model *Student) error {
	if err := tx.Create(model).Error; err != nil {
		return err
	}
	return nil
}

func DeleteStudent(id int64) error {
	if err := WriteDB().Where("studentId = ?", id).Delete(&Student{}).Error; err != nil {
		return err
	}
	return nil
}
