package util

/*utils 包内所有结构体在这里定义，解决循环引用问题*/

// token换取用户信息josn
type ApiStudent struct {
	Code int            `json:"code"`
	Data ApiStudentData `json:"data"`
}

type ApiStudentSlice struct {
	Code int              `json:"code"`
	Data []ApiStudentData `json:"data"`
}
type ApiStudentData struct {
	StudentId int64  `json:"studentId"`
	UserId    int64  `json:"userId"`
	Mobile    string `json:"mobile"`
	UserName  string `json:"username"`
	Token     string `json:"token"`
}

// 获取配置json
type ConfigInfo struct {
	Code int        `json:"code"`
	Data ConfigData `json:"data"`
}
type ConfigData struct {
	User []string `json:"user"`
	Sms  []string `json:"sms"`
}

//用户中心注册实体
type RegisterStudent struct {
	UnionId      string `json:"unionid"`
	Nickname     string `json:"nickname"`
	HeadPhoto    string `json:"headPhoto"`
	Sex          string `json:"sex"`
	Province     string `json:"province"`
	City         string `json:"city"`
	Country      string `json:"country"`
	ClientSource string `json:"clientSource"`
	Version      string `json:"version"`
	Channel      string `json:"channel"`
}

// 获取token返回信息
type ApiToken struct {
	Code int       `json:"code"`
	Data TokenData `json:"data"`
}

type TokenData struct {
	Token string `json:"token"`
}

// 获取api返回码
type ApiCode struct {
	Code int `json:"code"`
}
