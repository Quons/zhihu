package models

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"time"
)

type Lessonsection struct {
	SectionID int64  `gorm:"primary_key;column:sectionId"`
	LessonID  int64  `gorm:"column:lessonId"`
	Lesson    Lesson `gorm:"ForeignKey:lessonId"`
	//关联查询，以ForeignKey作为外键 从Course表中以AssociationForeignKey 为references进行查询
	Course             Course    `gorm:"ForeignKey:courseId"`
	CourseID           int64     `gorm:"column:courseId"`
	SectionNum         int       `gorm:"column:sectionNum"`
	Brief              string    `gorm:"column:brief"`
	SectionName        string    `gorm:"column:sectionName"`
	VideoURL           string    `gorm:"column:videoUrl"`
	VideoTimes         int       `gorm:"column:videoTimes"`
	Courseware         string    `gorm:"column:courseware"`
	ParentID           int64     `gorm:"column:parentId"`
	AddTime            time.Time `gorm:"column:addTime"`
	AnswerFile         string    `gorm:"column:answerFile"`
	WorksCover         string    `gorm:"column:worksCover"`
	CircleShareTitle   string    `gorm:"column:circleShareTitle"`
	CircleShareContent string    `gorm:"column:circleShareContent"`
	FriendShareTitle   string    `gorm:"column:friendShareTitle"`
	FriendShareContent string    `gorm:"column:friendShareContent"`
	IsChallenge        int8      `gorm:"column:isChallenge"`
	IsClockIn          int8      `gorm:"column:isClockIn"`
	SectionImage       string    `gorm:"column:sectionImage"`
}

func ExistLessonSectionByID(id int64) (bool, error) {
	var lessonSection Lessonsection
	err := readDB().Select("sectionId").Where("sectionId = ?", id).First(&lessonSection).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if lessonSection.SectionID > 0 {
		return true, nil
	}

	return false, nil
}

func GetLessonSectionTotal(maps interface{}) (int, error) {
	var count int
	if err := readDB().Model(&Lessonsection{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func GetLessonSections(pageNum int, pageSize int, maps interface{}) ([]Lessonsection, error) {
	var lessonSections []Lessonsection
	err := readDB().Order("sectionId desc").Offset(pageNum).Limit(pageSize).Where(maps).Find(&lessonSections).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return lessonSections, err
	}
	return lessonSections, nil
}

func GetLessonSection(id int64) (*Lessonsection, error) {
	lessonSection := Lessonsection{SectionID: id}
	err := readDB().First(&lessonSection).Related(&lessonSection.Course).Related(&lessonSection.Lesson).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logrus.Errorf("%+v", err)
		return nil, err
	}
	return &lessonSection, nil
}

func AddOrUpdateLessonSection(lessonSection *Lessonsection) error {
	if err := WriteDB().Save(lessonSection).Error; err != nil {
		return err
	}
	return nil
}

func AddLessonSection(lessonSection *Lessonsection) error {
	return AddLessonSectionTrans(WriteDB(), lessonSection)
}

func AddLessonSectionTrans(tx *gorm.DB, lessonSection *Lessonsection) error {
	if err := tx.Create(lessonSection).Error; err != nil {
		return err
	}
	return nil
}

func DeleteLessonSection(id int64) error {
	if err := WriteDB().Where("sectionId = ?", id).Delete(Lessonsection{}).Error; err != nil {
		return err
	}
	return nil
}
