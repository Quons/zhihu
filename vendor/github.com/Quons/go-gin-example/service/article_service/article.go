package article_service

import (
	"fmt"
	"github.com/Quons/go-gin-example/models"
	. "github.com/Quons/go-gin-example/pkg/e"
	"github.com/Quons/go-gin-example/pkg/gredis"
	"github.com/sirupsen/logrus"
)

type Article struct {
	ID            int
	TagID         int
	Title         string
	Desc          string
	Content       string
	CoverImageUrl string
	State         int
	CreatedBy     string
	ModifiedBy    string

	PageNum  int
	PageSize int
}

func (a *Article) Add() error {
	article := map[string]interface{}{
		"tag_id":          a.TagID,
		"title":           a.Title,
		"desc":            a.Desc,
		"content":         a.Content,
		"created_by":      a.CreatedBy,
		"cover_image_url": a.CoverImageUrl,
		"state":           a.State,
	}

	if err := models.AddArticle(article); err != nil {
		return err
	}

	return nil
}

func (a *Article) AddArticleAndTag() {
	logrus.WithFields(logrus.Fields{
		"name": "quon",
		"age":  25,
	}).Info("start transaction！！")
	tx := models.WriteDB().Begin()
	defer func() {
		if s := recover(); s != nil {
			logrus.Errorf("%v", s)
			tx.Rollback()
		}
	}()
	//添加文章
	article := map[string]interface{}{
		"tag_id":          a.TagID,
		"title":           a.Title,
		"desc":            a.Desc,
		"content":         a.Content,
		"created_by":      a.CreatedBy,
		"cover_image_url": a.CoverImageUrl,
		"state":           a.State,
	}
	if err := models.AddArticleTrans(tx, article); err != nil {
		tx.Rollback()
		fmt.Printf(err.Error())
		return
	}
	//添加标签
	if err := models.AddTagTrans(tx, "testTag", 1, "quons"); err != nil {
		tx.Rollback()
		fmt.Printf(err.Error())
		return
	}
	if err := tx.Commit().Error; err != nil {
		fmt.Printf(err.Error())
		logrus.Errorf("%+v", err)
	}

}

func (a *Article) Edit() error {
	return models.EditArticle(a.ID, map[string]interface{}{
		"tag_id":          a.TagID,
		"title":           a.Title,
		"desc":            a.Desc,
		"content":         a.Content,
		"cover_image_url": a.CoverImageUrl,
		"state":           a.State,
		"modified_by":     a.ModifiedBy,
	})
}

func (a *Article) Get() (*models.Article, error) {
	var article *models.Article
	if gredis.Get(article, CACHE_ARTICLE, a.ID) {
		return article, nil
	}

	article, err := models.GetArticle(a.ID)
	if err != nil {
		return nil, err
	}

	gredis.Set(article, CACHE_ARTICLE_TIME, CACHE_TAG, article.ID)
	return article, nil
}

func (a *Article) GetAll() (models.Articles, error) {
	var articles models.Articles

	if gredis.Get(&articles, CACHE_ARTICLE_LIST, a.TagID, a.State, a.PageNum, a.PageNum) {
		return articles, nil
	}

	articles, err := models.GetArticles(a.PageNum, a.PageSize, a.getMaps())
	if err != nil {
		return nil, err
	}

	gredis.Set(articles, CACHE_ARTICLE_TIME, CACHE_ARTICLE_LIST, a.TagID, a.State, a.PageNum, a.PageNum)
	return articles, nil
}

func (a *Article) Delete() error {
	return models.DeleteArticle(a.ID)
}

func (a *Article) ExistByID() (bool, error) {
	return models.ExistArticleByID(a.ID)
}

func (a *Article) Count() (int, error) {
	return models.GetArticleTotal(a.getMaps())
}

func (a *Article) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0
	if a.State != -1 {
		maps["state"] = a.State
	}
	if a.TagID != -1 {
		maps["tag_id"] = a.TagID
	}

	return maps
}
