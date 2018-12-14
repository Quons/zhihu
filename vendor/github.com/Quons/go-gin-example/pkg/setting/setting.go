package setting

import (
	"log"
	"time"

	"github.com/Quons/go-gin-example/pkg/file"
	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

type App struct {
	JwtSecret string
	PageSize  int
	PrefixUrl string

	RuntimeRootPath string

	ImageSavePath  string
	ImageMaxSize   int
	ImageAllowExts []string

	ExportSavePath string
	QrCodeSavePath string
	FontSavePath   string

	LogLevel string
}

var AppSetting = &App{}

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var ServerSetting = &Server{}

type Database struct {
	Type        string
	Name        string
	TablePrefix string
	WUser       string
	WPassword   string
	WHost       string
	RUser       string
	RPassword   string
	RHost       string
}

var DatabaseSetting = &Database{}

type Redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
	Prefix      string
}

var RedisSetting = &Redis{}

type UserCenterSetting struct {
	ConfigUrl1             string `ini:"getConfigUrl1"`
	ConfigUrl2             string `ini:"getConfigUrl2"`
	GetStudentByTokenUrl   string `ini:"getStudentByTokenUrl"`
	GetTokenByStudentIdUrl string `ini:"getTokenByStudentIdUrl"`
	WeChatCreateStudentUrl string `ini:"weChatCreateStudentUrl"`
	WeChatLoginUrl         string `ini:"weChatLoginUrl"`
}

var UserCenter = &UserCenterSetting{}

var cfg *ini.File

func Setup(runmode string) {
	configFile := ""
	switch runmode {
	case "dev":
		configFile = "dev.ini"
	case "test":
		configFile = "test.ini"
	case "pre":
		configFile = "pre.ini"
	case "prod":
		configFile = "prod.ini"
	default:
		logrus.Fatal("invalid runmode,must be one of [dev,test,pre,prod]!")
	}
	//获取绝对路径
	workPath, err := file.GetExecPath()
	if err != nil {
		logrus.Fatalf("get execPath error:%+v", err)
	}
	cfg, err = ini.Load(filepath.Join(workPath, "conf", configFile))
	if err != nil {
		//根据goPath定位配置文件
		goPath := os.Getenv("GOPATH")
		cfg, err = ini.Load(filepath.Join(goPath, "src/github.com/Quons/go-gin-example", "conf", configFile))
		if err != nil {
			logrus.Fatalf("get execPath error:%+v", err)
		}
	}

	mapTo("app", AppSetting)
	mapTo("server", ServerSetting)
	mapTo("database", DatabaseSetting)
	mapTo("redis", RedisSetting)
	mapTo("userCenter", UserCenter)

	AppSetting.ImageMaxSize = AppSetting.ImageMaxSize * 1024 * 1024
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.ReadTimeout * time.Second
	RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second
}

func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo RedisSetting err: %v", err)
	}
}
