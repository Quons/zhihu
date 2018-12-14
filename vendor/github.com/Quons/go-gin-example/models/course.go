package models

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"time"
)

type Course struct {
	CourseID   int64  `gorm:"primary_key;column:courseId"`
	CourseName string `gorm:"column:courseName"`
	//一对多，指定的是多的一方的外键以及对于一的一方的reference
	Lessons            []Lesson        `gorm:"ForeignKey:lessonId"`
	LessonSections     []Lessonsection `gorm:"ForeignKey:courseId"`
	CourseImage        string          `gorm:"column:courseImage"`
	Stage              string          `gorm:"column:stage"`
	LessonNum          int             `gorm:"column:lessonNum"`
	OriginalPrice      float64         `gorm:"column:originalPrice"`
	Price              float64         `gorm:"column:price"`
	LimitNum           int             `gorm:"column:limitNum"`
	BuyNum             int             `gorm:"column:buyNum"`
	MallProductID      int             `gorm:"column:mallProductId"`
	Status             int8            `gorm:"column:status"`
	TeacherID          int64           `gorm:"column:teacherId"`
	AssistantTeacherID int64           `gorm:"column:assistantTeacherId"`
	AddTime            time.Time       `gorm:"column:addTime"`
	StartTime          time.Time       `gorm:"column:startTime"`
	EndTime            time.Time       `gorm:"column:endTime"`
	CourseType         int8            `gorm:"column:courseType"`
	IsExperienceCourse int8            `gorm:"column:isExperienceCourse"`
	ListImage          string          `gorm:"column:listImage"`
	ScratchFile        string          `gorm:"column:scratchFile"`
	PrizeName          string          `gorm:"column:prizeName"`
	PrizeImage         string          `gorm:"column:prizeImage"`
	PrizeDetail        string          `gorm:"column:prizeDetail"`
	NeedBuyCourse      int8            `gorm:"column:needBuyCourse"`
	Summary            string          `gorm:"column:summary"`
	AnswerFile         string          `gorm:"column:answerFile"`
	OpenSectionID      int64           `gorm:"column:openSectionId"`
}

/*// 设置User的表名为`profiles`
func (Course) TableName() string {
	return "cb_course"
}*/

func ExistCourseByID(id int64) (bool, error) {
	var course Course
	err := readDB().Select("courseId").Where("courseId = ?", id).First(&course).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if course.CourseID > 0 {
		return true, nil
	}

	return false, nil
}

func GetCourseTotal(maps interface{}) (int, error) {
	var count int
	if err := readDB().Model(&Course{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func GetCourses(pageNum int, pageSize int, maps interface{}) ([]Course, error) {
	var courses []Course
	err := readDB().Order("courseId desc").Offset(pageNum).Limit(pageSize).Where(maps).Find(&courses).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return courses, err
	}
	return courses, nil
}

func GetCourse(id int64) (*Course, error) {
	var course Course
	//err := readDB().Preload("LessonSections", func(db *gorm.DB) *gorm.DB { return db.Where("sectionNum > 4") }).Preload("Lessons").First(&course).Error
	err := readDB().Where("courseId=?", id).First(&course).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logrus.Errorf("%+v", err)
		return nil, err
	}
	return &course, nil
}

//id为0时新增，否则更新
func AddOrUpdateCourse(course *Course) error {
	if err := WriteDB().Save(course).Error; err != nil {
		return err
	}
	return nil
}

func AddCourse(course *Course) error {
	return AddCourseTrans(WriteDB(), course)
}

func AddCourseTrans(tx *gorm.DB, course *Course) error {
	if err := tx.Create(course).Error; err != nil {
		return err
	}
	return nil
}

func DeleteCourse(id int64) error {
	if err := WriteDB().Where("courseId = ?", id).Delete(Course{}).Error; err != nil {
		return err
	}
	return nil
}
