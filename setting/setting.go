package setting

/*
	在应用程序的运行生命周期中, 最直接的关系之一就是应用的配置读取和更新;
	它的一举一动都有可能影响应用程序的改变, 其分别包含如下行为:
	- 在启动时:
		可以进行一些基础应用属性、连接第三方实例（MySQL、NoSQL）等等的初始化
		行为;
	- 在运行中:
		可以监听文件或其他存储载体的变更来实现热更新配置的效果, 例如: 在发现
		有变更的话, 就对原有配置值进行修改以此达到相关联的一个效果;
		如果更深入业务的使用, 还可以通过配置的热更新, 达到功能灰度(TODO)的效果,
		也是一个比较常见的场景;

	Viper 是适用于 Go 应用程序的完整配置解决方案, 是目前 Go 语言中比较流行
	的文件配置解决方案, 支持处理各种不同类型的配置需求和配置格式; TODO: 熟悉
*/

import "github.com/spf13/viper"

type Setting struct {
	vp *viper.Viper
}

func NewSetting(confFn string) (*Setting, error) {
	vp := viper.New()
	vp.SetConfigName("config")
	vp.AddConfigPath("configs/")
	vp.SetConfigType("yaml")
	if confFn != "" {
		vp.SetConfigFile(confFn)
	}

	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return &Setting{vp}, nil
}
