package util

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

func HttpGet(getUrl string) (string, error) {
	resp, err := http.Get(getUrl)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), err
}

func HttpPost(postUrl string, postData map[string]string) (string, error) {
	postValue := url.Values{}
	for key, value := range postData {
		postValue.Set(key, value)
	}
	data := postValue.Encode()
	resp, err := http.Post(postUrl,
		"application/x-www-form-urlencoded",
		strings.NewReader(data))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), err
}

//HTTP 包中POST 方法
func HttpPost2(my_url string, postData string) (string, error) {
	//HTTP POST请求
	resp, err := http.Post(my_url,
		"application/json;charset=utf-8",
		strings.NewReader(postData))
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), err
}

//结构体所有字段必须为string类型，model不能是指针类型
func ToPostDate(model interface{}) (map[string]string, error) {
	postData := make(map[string]string)
	t := reflect.TypeOf(model)
	v := reflect.ValueOf(model)
	for k := 0; k < t.NumField(); k++ {
		f := v.Field(k)
		switch f.Kind() {
		case reflect.String:
			postData[t.Field(k).Tag.Get("json")] = f.String()
		default:
			return postData, errors.New("invalid post data type,must be string")
		}
	}
	return postData, nil
}
