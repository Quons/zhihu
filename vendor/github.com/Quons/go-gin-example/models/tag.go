package models

import (
	"github.com/jinzhu/gorm"
)

type Tag struct {
	Model
	Articles   []Article
	Name       string `gorm:"column:name"`
	CreatedBy  string `gorm:"column:created_by"`
	ModifiedBy string `gorm:"column:modified_by"`
	State      int    `gorm:"column:state"`
}

func ExistTagByName(name string) (bool, error) {
	var tag Tag
	err := readDB().Select("id").Where("name = ? AND deleted_on = ? ", name, 0).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if tag.ID > 0 {
		return true, nil
	}
	return false, nil
}

func AddTag(name string, state int, createdBy string) error {
	return AddTagTrans(WriteDB(), name, state, createdBy)
}

func AddTagTrans(tx *gorm.DB, name string, state int, createdBy string) error {
	tag := Tag{
		Name:      name,
		State:     state,
		CreatedBy: createdBy,
	}
	if err := tx.Create(&tag).Error; err != nil {
		return err
	}

	return nil
}

func GetTags(pageNum int, pageSize int, maps interface{}) ([]Tag, error) {
	var tags []Tag
	readDb := readDB()
	if pageSize >= 0 && pageNum > 0 {
		readDb = readDb.Offset(pageNum).Limit(pageSize)
	}

	err := readDb.Order("id").Preload("Articles").Where(maps).Find(&tags).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return tags, nil
}

func GetTagTotal(maps interface{}) (int, error) {
	var count int
	if err := readDB().Model(&Tag{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func ExistTagByID(id int) (bool, error) {
	var tag Tag
	err := readDB().Select("id").Where("id = ? AND deleted_on = ? ", id, 0).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if tag.ID > 0 {
		return true, nil
	}

	return false, nil
}

func DeleteTag(id int) error {
	if err := WriteDB().Where("id = ?", id).Delete(&Tag{}).Error; err != nil {
		return err
	}

	return nil
}

func EditTag(id int, data interface{}) error {
	if err := WriteDB().Model(&Tag{}).Where("id = ? AND deleted_on = ? ", id, 0).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

func CleanAllTag() (bool, error) {
	if err := WriteDB().Unscoped().Where("deleted_on != ? ", 0).Delete(&Tag{}).Error; err != nil {
		return false, err
	}

	return true, nil
}
