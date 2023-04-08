package setting

import "time"

type ServerSettingS struct {
	RunMode        string
	HttpPort       string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	ContextTimeout time.Duration
}

type AppSettingS struct {
	DefaultPageSize      int
	MintLimit            int
	MaxPageSize          int
	LogSavePath          string
	LogFileName          string
	LogFileExt           string
	NftMeetLimit         int
	MixEventHours        int
	MintCountLimit       int
	UploadSavePath       string
	UploadServerUrl      string
	UploadImageMaxSize   int
	UploadImageAllowExts []string
}

type DatabaseSettingS struct {
	DBType       string
	UserName     string
	Password     string
	Host         string
	DBName       string
	TablePrefix  string
	Charset      string
	ParseTime    bool
	MaxIdleConns int
	MaxOpenConns int
}

func (s *Setting) ReadSection(k string, v interface{}) error {
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}

	return nil
}

type JWTSettingS struct {
	Secret string
}

type EmailSettingS struct {
	Host     string
	Port     int
	UserName string
	Password string
	IsSSL    bool
	From     string
	To       []string
}

type APISettingS struct {
	ScPre  string
	ScProd string
	AmPre  string
	AmProd string
	WxProd string
}

type WxAPISettingS struct {
	GrantType string
	AppId     string
	Secret    string
}

type AliAccountSettingS struct {
	ApiId     string
	ApiSecret string
	Region    string
}

type StableAPISettingS struct {
	API string
	Key string
}
