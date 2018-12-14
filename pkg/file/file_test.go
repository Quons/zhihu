package file

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetExecPath(t *testing.T) {
	dirName, err := GetDirName()
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, dirName, "go-gin-example")
}
