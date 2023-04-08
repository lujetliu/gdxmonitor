package global

import (
	"flag"
	"gdxmonitor/logger"
	"gdxmonitor/setting"
	"io"
	"log"
	"os"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

var TIME_FORMAT = "2006-01-02 15:04:05"

var (
	ServerSetting   *setting.ServerSettingS
	AppSetting      *setting.AppSettingS
	DatabaseSetting *setting.DatabaseSettingS
	Logger          *logger.Logger
)

// TODO: 应该在main文件中实现
var confParam = flag.String("f", "configs/config.yaml", "config.yml's path")

func init() {
	flag.Parse()
	_, err := os.Stat(*confParam)
	if os.IsNotExist(err) {
		panic(err)
	}

	err = setupSetting(*confParam)
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}

	err = setupLogger()
	if err != nil {
		log.Fatalf("init.setupLogger err: %v", err)
	}
}

func setupLogger() error {
	writer := &lumberjack.Logger{
		Filename: AppSetting.LogSavePath + "/" +
			AppSetting.LogFileName + AppSetting.LogFileExt,
		MaxSize:   300, // 日志文件允许的最大占用空间为 300 MB
		MaxAge:    100, // 日志文件最大生存周期为100天 TODO: 如何监控文件过期的
		LocalTime: true,
	}
	multiWriter := io.MultiWriter(os.Stdout, writer)
	Logger = logger.NewLogger(multiWriter, "", log.LstdFlags)

	return nil
}

func setupSetting(confFn string) error {
	setting, err := setting.NewSetting(confFn)
	if err != nil {
		return err
	}
	err = setting.ReadSection("Server", &ServerSetting)
	if err != nil {
		return err
	}
	ServerSetting.ReadTimeout *= time.Second
	ServerSetting.WriteTimeout *= time.Second

	err = setting.ReadSection("App", &AppSetting)
	if err != nil {
		return err
	}

	err = setting.ReadSection("Database", &DatabaseSetting)
	if err != nil {
		return err
	}
	return nil
}
