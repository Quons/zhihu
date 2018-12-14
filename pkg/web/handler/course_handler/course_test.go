package course_handler

import (
	"github.com/Quons/go-gin-example/models"
	"github.com/Quons/go-gin-example/pkg/setting"
	"testing"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/assert"
)

func init() {
	setting.Setup("dev")
	models.Setup()
}

type MockedObj struct {
	mock.Mock
}

func (m *MockedObj) Speek(word string) (s string) {
	arg := m.Called(word)
	return arg.String(0)
}

func TestAddArticleAndTag(t *testing.T) {
	testObj := new(MockedObj)
	testObj.On("Speek",mock.Anything).Return("hhhhh")
	s:=GetTest(testObj)
	t.Log(s)
	assert.Equal(t,"hhhhh",s)
	testObj.AssertCalled(t,"Speek","hello")
}
