package v1

import (
	"strconv"
	"testing"
	"github.com/pkg/errors"
	"fmt"
)

func TestStr(t *testing.T) {
	s, err := strconv.ParseInt("d", 10, 0)
	fmt.Println(err)
	err=errors.WithStack(err)
	if err != nil {
		t.Errorf("%v",err)
		return
	}
	t.Log(s)
}
