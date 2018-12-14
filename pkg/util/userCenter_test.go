package util

import (
	"github.com/Quons/go-gin-example/pkg/gredis"
	"github.com/Quons/go-gin-example/pkg/setting"
	"testing"
)

func init() {
	setting.Setup("dev")
	gredis.Setup()
}

func TestGetApiStudentByUnionId(t *testing.T) {
	apiStudent, err := GetApiStudentByUnionId("oMqARszvlqjFWeg6FDhVaXqfZZck")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(apiStudent)
}

func TestCreateStudentFromUserCenter(t *testing.T) {
	rs := RegisterStudent{UnionId: "oMqARszvlqjFWeg6FDhVaXqfZZck", Nickname: "Quon", HeadPhoto: "https://appd.knowbox.cn/codebox/hp/data/go/src/bianchengapi/tmp/8bd63d20827bffb1dd724fe401f3deeb.png",
		Sex: "male", Province: "Beijing", City: "Beijing", Country: "China", ClientSource: "wx", Version: "1.0", Channel: "wx"}
	apiStudent, err := CreateStudentFromUserCenter(rs)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(apiStudent)
}
