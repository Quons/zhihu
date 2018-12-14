package util

import (
	"encoding/json"
	"github.com/Quons/go-gin-example/pkg/e"
	"github.com/Quons/go-gin-example/pkg/gredis"
	"github.com/Quons/go-gin-example/pkg/setting"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

// 获取配置中心的配置
func GetUserCenterConfig(forceUpdate bool) ConfigInfo {
	var configInfo ConfigInfo
	flag := 0
	gredis.Get(&flag, e.CACHE_USER_CENTER_CONF_FLAG)
	if !forceUpdate && flag == 1 {
		if gredis.Get(&configInfo, e.CACHE_USER_CENTER_CONF) {
			return configInfo
		}
	}
	//flag 过期，重新获取
	gredis.Set(1, e.CACHE_USER_CENTER_CONF_FLAG_TIME, e.CACHE_USER_CENTER_CONF_FLAG)
	config, httpErr := HttpGet(setting.UserCenter.ConfigUrl1)
	jsonErr := json.Unmarshal([]byte(config), &configInfo)
	if httpErr != nil || jsonErr != nil || configInfo.Code != 99999 {
		log.WithFields(log.Fields{"httpErr": httpErr, "jsonError": jsonErr, "configInfo": configInfo}).Error("get user config err")
		//发生错误，重试
		config, httpErr = HttpGet(setting.UserCenter.ConfigUrl2)
		jsonErr = json.Unmarshal([]byte(config), &configInfo)
		if httpErr != nil || jsonErr != nil || configInfo.Code != 99999 {
			log.WithFields(log.Fields{"httpErr": httpErr, "jsonErr": jsonErr, "configInfo": configInfo}).Error("get user center config err")
			return configInfo
		}
	}
	gredis.Set(configInfo, e.CACHE_USER_CENTER_CONF_TIME+86400, e.CACHE_USER_CENTER_CONF)
	return configInfo
}

// 得到用户中心的接口ip
func getUserCenterIp(token string) (string, error) {
	var userCenterIp string
	if gredis.Get(&userCenterIp, e.CACHE_USER_CENTER_IP, token) || userCenterIp != "" {
		return userCenterIp, nil
	}

	var tokenHash = Crc32IEEE(token)
	var configInfo ConfigInfo
	configInfo = GetUserCenterConfig(false)
	userCenterIp = configInfo.Data.User[tokenHash%uint32(len(configInfo.Data.User))]
	if userCenterIp == "" {
		return "", errors.New("empty userCenterIp")
	}
	gredis.Set(userCenterIp, e.CACHE_USER_CENTER_IP_TIME, e.CACHE_USER_CENTER_IP)
	return userCenterIp, nil
}

/*通过token从用户中心获取学生信息*/

func GetStudentFromUserCenter(token string) (ApiStudent, error) {
	var apiStudent ApiStudent
	userCenterIp, err := getUserCenterIp(token)
	if err != nil {
		return apiStudent, err
	}
	var postData = make(map[string]string)
	postData["token"] = token
	reqUrl := strings.Replace(setting.UserCenter.GetStudentByTokenUrl, "{ip}", userCenterIp, 1)
	studentInfo, err := HttpPost(reqUrl, postData)
	if err != nil {
		return apiStudent, err
	}
	if studentInfo == "" {
		return apiStudent, errors.New("fail to get userinfo from user center,token:" + token)
	}
	var apiCode ApiCode
	err = json.Unmarshal([]byte(studentInfo), &apiCode)
	if err != nil {
		log.WithField("studentInfo", studentInfo).Error(err)
		return apiStudent, err
	}

	if apiCode.Code == 30003 {
		log.WithField("studentInfo", studentInfo).Error()
		return apiStudent, errors.New("30003")
	}

	if apiCode.Code == 99999 {
		err = json.Unmarshal([]byte(studentInfo), &apiStudent)
		if err != nil {
			log.WithField("studentInfo", studentInfo).Error(err)
			return apiStudent, err
		}
		return apiStudent, nil
	}
	return apiStudent, errors.New(studentInfo)
}

/*通过studentId 获取用户中心token*/
func GetApiTokenByStudentId(studentId int64) (ApiToken, error) {
	var apiToken ApiToken
	/*获取用户中心ip*/
	userCenterIp, err := getUserCenterIp(strconv.FormatInt(studentId, 10))
	if err != nil {
		log.WithField("studentId", studentId).Error(err)
		return apiToken, err
	}
	var postData = make(map[string]string)
	postData["studentId"] = strconv.FormatInt(studentId, 10)
	reqUrl := strings.Replace(setting.UserCenter.GetTokenByStudentIdUrl, "{ip}", userCenterIp, 1)
	tokenInfo, err := HttpPost(reqUrl, postData)
	if err != nil {
		log.WithFields(log.Fields{"reqUrl": reqUrl, "postData": postData}).Error(err)
		return apiToken, err
	}
	if tokenInfo == "" {
		log.WithField("studentId", studentId).Error("get no token")
		return apiToken, errors.New("empty token response")
	}
	var apiCode ApiCode
	err = json.Unmarshal([]byte(tokenInfo), &apiCode)
	if err != nil {
		log.WithField("tokenInfo", tokenInfo).Error(err)
		return apiToken, err
	}
	if apiCode.Code != 99999 {
		log.WithFields(log.Fields{"studentId": studentId, "tokenInfo": tokenInfo}).Error()
		return apiToken, errors.New(tokenInfo)
	}
	err = json.Unmarshal([]byte(tokenInfo), &apiToken)
	if err != nil {
		log.WithField("tokenInfo", tokenInfo).Error(err)
		return apiToken, err
	}
	if apiToken.Data.Token == "" {
		return apiToken, errors.New("empty token,studentId:" + strconv.FormatInt(studentId, 10))
	}
	return apiToken, nil
}

/*在用户中心创建学生用户并返回*/

func CreateStudentFromUserCenter(registerStudent RegisterStudent) (ApiStudent, error) {
	var apiStudent ApiStudent
	userCenterIp, err := getUserCenterIp(registerStudent.UnionId)
	if err != nil {
		return apiStudent, err
	}
	var postData2 = make(map[string]string)
	postData2["unionid"] = registerStudent.UnionId
	postData2["nickname"] = registerStudent.Nickname
	postData2["headPhoto"] = registerStudent.HeadPhoto
	postData2["sex"] = registerStudent.Sex
	postData2["province"] = registerStudent.Province
	postData2["city"] = registerStudent.City
	postData2["country"] = registerStudent.Country
	postData2["clientSource"] = registerStudent.ClientSource
	postData2["version"] = registerStudent.Version
	postData2["channel "] = registerStudent.Channel
	reqUrl := strings.Replace(setting.UserCenter.WeChatCreateStudentUrl, "{ip}", userCenterIp, 1)
	tokenInfo, err := HttpPost(reqUrl, postData2)
	if err != nil {
		log.WithFields(log.Fields{"userCententIp": userCenterIp, "postData": postData2}).Error(err)
		return apiStudent, errors.New("wechatCreateStudentUrl error")
	}
	if tokenInfo == "" {
		return apiStudent, errors.New("wechatCreateStudentUrl error:empty response")
	}
	var apiCode ApiCode
	err = json.Unmarshal([]byte(tokenInfo), &apiCode)
	if err != nil {
		log.WithField("tokenInfo", tokenInfo).Error(err)
		return apiStudent, err
	}
	if apiCode.Code == 99999 {
		//返回码正确，进行正常解析
		err = json.Unmarshal([]byte(tokenInfo), &apiStudent)
		if err != nil {
			log.WithFields(log.Fields{"tokenInfo": tokenInfo, "reqUrl": reqUrl, "postData": postData2}).Error(err)
			return apiStudent, err
		}
		return apiStudent, nil
	} else if apiCode.Code == 30000 || apiCode.Code == 30004 {
		//说明用户已注册或已绑定账户 ，查询用户token信息，
		apiStudent, err = GetApiStudentByUnionId(registerStudent.UnionId)
		if err != nil {
			log.WithField("unionId", registerStudent.UnionId).Error(err)
			return apiStudent, err
		}
		return apiStudent, nil
	} else {
		log.WithField("tokenInfo", tokenInfo).Error()
		return apiStudent, errors.New("get api student error")
	}
}

/*根据unionId获取用户中心学生信息*/

func GetApiStudentByUnionId(unionId string) (ApiStudent, error) {
	var apiStudent ApiStudent
	userCenterIp, err := getUserCenterIp(unionId)
	if err != nil {
		log.WithField("unionId", unionId).Error(err)
		return apiStudent, err
	}
	var postData3 = make(map[string]string)
	postData3["unionid"] = unionId
	reqUrl := strings.Replace(setting.UserCenter.WeChatLoginUrl, "{ip}", userCenterIp, 1)
	tokenInfo, err := HttpPost(reqUrl, postData3)
	if err != nil {
		log.WithFields(log.Fields{"reqUrl": reqUrl, "postData": postData3}).Error(err)
		return apiStudent, errors.New("wechatLoginUrl error")
	}
	if tokenInfo == "" {
		log.WithField("unionId", unionId).Error("empty token info")
		return apiStudent, errors.New("wechatLoginUrl error")
	}
	var apiCode ApiCode
	err = json.Unmarshal([]byte(tokenInfo), &apiCode)
	if err != nil {
		log.WithField("tokenInfo", tokenInfo).Error(err)
		return apiStudent, err
	}
	if apiCode.Code != 99999 {
		log.WithFields(log.Fields{"tokenInfo": tokenInfo, "unionId": unionId}).Error("get api student err")
		return apiStudent, errors.New(tokenInfo)
	}
	//返回正常，解析
	err = json.Unmarshal([]byte(tokenInfo), &apiStudent)
	if err != nil {
		log.WithField("tokenInfo", tokenInfo).Error(err)
		return apiStudent, err
	}
	return apiStudent, nil
}
