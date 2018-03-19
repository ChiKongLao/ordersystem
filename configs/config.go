package configs

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Environment string      `yaml:"Environment"`
	Mqtt        MqttConfig  `yaml:"Mqtt"`
	MySQL       MySQLConfig `yaml:"MySQL"`
	Redis       RedisConfig `yaml:"Redis"`
}
type MySQLConfig struct {
	Host     string `yaml:"Host"`
	UserName string `yaml:"UserName"`
	Password string `yaml:"Password"`
	Name     string `yaml:"Name"`
}

type RedisConfig struct {
	Host string `yaml:"Host"`
}

type MqttConfig struct {
	Host string `yaml:"Host"`
	Open bool   `yaml:"Open"`
}

var mConfig Config

func InitApplicationConfig() {

	defer func() {
		if err := recover(); err != nil {
			logrus.Warn("加载配置文件失败: ", err)
			mConfig = Config{
				Environment: "test",
				Mqtt: MqttConfig{
					Host: "tcp://106.38.205.89:1883",
					Open: true,
				},
				MySQL: MySQLConfig{
					Host:     "localhost:3306",
					UserName: "root",
					Password: "",
					Name:     "order_system",
				},
				Redis: RedisConfig{
					Host: "localhost:6379",
				},
			}

		}
	}()

	file := "./configs/test.yml"
	//file := "./configs/prod.yml"

	data, err := ioutil.ReadFile(file)

	if err != nil { // 不为空时使用默认
		panic(err)
	}
	err = yaml.Unmarshal(data, &mConfig)
	if err != nil {
		panic(err)
	}

	logrus.Infoln("加载配置文件成功: ", mConfig)

}

func GetConfig() *Config {
	return &mConfig
}

// 是否为测试环境
func (config *Config) IsTest() bool {
	return config.Environment == "test"

}

// 是否为预生产环境
func (config *Config) IsPreProd() bool {
	return config.Environment == "pre_prod"
}

// 是否为生产环境
func (config *Config) IsProd() bool {
	return config.Environment == "prod"
}

func (config *Config) GetEnvironmentName() string {
	switch config.Environment {
	case "test":
		return "测试环境"
	case "pre_prod":
		return "预生产环境"
	case "prod":
		return "生产环境"
	}
	return "未知"
}
