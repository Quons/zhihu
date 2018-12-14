package models

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"sort"
)

//Article 文章
type Article struct {
	Model

	TagID         int    `gorm:"column:tag_id"`
	Title         string `gorm:"column:title"`
	Tag           Tag
	Desc          string `gorm:"column:desc"`
	Content       string `gorm:"column:content"`
	CoverImageURL string `gorm:"column:cover_image_url"`
	CreatedBy     string `gorm:"column:created_by"`
	ModifiedBy    string `gorm:"column:modified_by"`
	State         int    `gorm:"column:state"`
}

type Articles []Article

func (model Articles) Len() int {
	return len(model)
}

func (model Articles) Swap(i, j int) {
	model[i], model[j] = model[j], model[i]
}

func (model Articles) Less(i, j int) bool { // 重写 Less() 方法， 从小到大排序
	return model[i].Tag.ID < model[j].Tag.ID
}

func ExistArticleByID(id int) (bool, error) {
	var article Article
	err := readDB().Select("id").Where("id = ? AND deleted_on = ? ", id, 0).First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if article.ID > 0 {
		return true, nil
	}

	return false, nil
}

func GetArticleTotal(maps interface{}) (int, error) {
	var count int
	if err := readDB().Model(&Article{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func GetArticles(pageNum int, pageSize int, maps interface{}) (Articles, error) {
	db := WriteDB().Begin()
	var articles Articles
	err := db.Preload("Tag", func(db *gorm.DB) *gorm.DB {
		return db.Order("id ASC")
	}).Order("tag_id desc").Offset(pageNum).Find(&articles).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	db.Commit()
	sort.Sort(articles)
	return articles, nil
}

func GetArticle(id int) (*Article, error) {
	var article Article
	err := readDB().Preload("Tag").Where("id = ? AND deleted_on = ? ", id, 0).First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logrus.Errorf("%+v", err)
		return nil, err
	}
	return &article, nil
}

func EditArticle(id int, data interface{}) error {
	if err := WriteDB().Model(&Article{}).Where("id = ? AND deleted_on = ? ", id, 0).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

func AddArticle(data map[string]interface{}) error {
	return AddArticleTrans(WriteDB(), data)
}

func AddArticleTrans(tx *gorm.DB, data map[string]interface{}) error {
	article := Article{
		TagID:         data["tag_id"].(int),
		Title:         data["title"].(string),
		Desc:          data["desc"].(string),
		Content:       data["content"].(string),
		CreatedBy:     data["created_by"].(string),
		State:         data["state"].(int),
		CoverImageURL: data["cover_image_url"].(string),
	}
	if err := tx.Create(&article).Error; err != nil {
		return err
	}
	return nil
}

func DeleteArticle(id int) error {
	if err := WriteDB().Where("id = ?", id).Delete(Article{}).Error; err != nil {
		return err
	}

	return nil
}

func CleanAllArticle() error {
	if err := WriteDB().Unscoped().Where("deleted_on != ? ", 0).Delete(&Article{}).Error; err != nil {
		return err
	}

	return nil
}
