package article_handler

import (
	"github.com/Quons/go-gin-example/models"
	"github.com/Quons/go-gin-example/pkg/setting"
	"testing"
)

func init() {
	setting.Setup("dev")
	models.Setup()
}

func TestAddArticleAndTag(t *testing.T) {
	a := &Article{TagID: 1, Title: "testArticle", Desc: "hiahia", Content: "testContent", CreatedBy: "quon", CoverImageUrl: "http"}
	a.AddArticleAndTag()
}
