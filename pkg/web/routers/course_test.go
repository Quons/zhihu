package routers

import (
	"zhihu/pkg/web/model"
	"zhihu/pkg/gredis"
	"zhihu/pkg/logging"
	"zhihu/pkg/setting"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func init() {
	setting.Setup("dev")
	logging.Setup()
	model.Setup()
	gredis.Setup()
}

func TestPingRoute(t *testing.T) {
	router := InitRouter()

	w := httptest.NewRecorder()

	postData := url.Values{}
	postData.Set("token", "feahvJWLZP88FddSbhv_1D1t5oep_paHjg-VRjT-buJMVgOpXvrqGjDhgs1mxKFP")
	postData.Set("courseId", "4")
	req, _ := http.NewRequest("POST", "/api/v1/getCourse?"+postData.Encode(), strings.NewReader(postData.Encode()))
	req.Header.Set("content-type", "application/x-www-form-urlencoded;charset=utf-8")

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	t.Log(w.Body.String())
	assert.Equal(t, true, strings.Contains(w.Body.String(), "99999"))
}
