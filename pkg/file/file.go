package file

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"mime/multipart"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func GetSize(f multipart.File) (int, error) {
	content, err := ioutil.ReadAll(f)

	return len(content), err
}

func GetExt(fileName string) string {
	return path.Ext(fileName)
}

func CheckNotExist(src string) bool {
	_, err := os.Stat(src)

	return os.IsNotExist(err)
}

func CheckPermission(src string) bool {
	_, err := os.Stat(src)

	return os.IsPermission(err)
}

func IsNotExistMkDir(src string) error {
	if notExist := CheckNotExist(src); notExist == true {
		if err := MkDir(src); err != nil {
			return err
		}
	}

	return nil
}

func MkDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func MustOpen(fileName, filePath string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err: %v", err)
	}

	src := dir + "/" + filePath
	perm := CheckPermission(src)
	if perm == true {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	err = IsNotExistMkDir(src)
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExistMkDir src: %s, err: %v", src, err)
	}

	f, err := Open(src+fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("Fail to OpenFile :%v", err)
	}

	return f, nil
}

/*得到当前执行文件的的路径*/
func GetExecPath() (string, error) {
	execFile, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path2, err := filepath.Abs(execFile)
	if err != nil {
		return "", err
	}
	rst := filepath.Dir(path2)
	return rst, nil
}

/*在项目目录下创建相对目录 make relative directory*/
func MkRdir(p string) (string, error) {
	rst, err := GetExecPath()
	if err != nil {
		return "", err
	}
	var isExistDir bool
	dstPath := filepath.Join(rst, p)
	_, err = os.Stat(dstPath)
	if err == nil {
		isExistDir = true
	}
	if os.IsNotExist(err) {
		isExistDir = false
	}
	if !isExistDir {
		os.Mkdir(dstPath, os.ModePerm)
	}
	return dstPath, nil
}
func GetDirName() (string, error) {
	/*获取可执行文件名称*/
	execFile, err := GetExecPath()
	if err != nil {
		logrus.Error(err.Error())
		return "", err
	}
	lastIndex := strings.LastIndex(execFile, "/")
	execFile = string([]rune(execFile)[lastIndex+1:])
	return execFile, nil
}
