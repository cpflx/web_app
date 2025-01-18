package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Conf 全局变量，用来保存程序的所有配置信息
var Conf = new(App)

type App struct {
	*AppConfig   `mapstructure:"app"`
	*LogConfig   `mapstructure:"log"`
	*MysqlConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}

type AppConfig struct {
	Name      string `mapstructure:"name"`
	Mode      string `mapstructure:"mode"`
	Port      int    `mapstructure:"port"`
	Version   string `mapstructure:"version"`
	StartTime string `mapstructure:"start_time"`
	MachineID int64  `mapstructure:"machine_id"`
}

type LogConfig struct {
	Filename   string `mapstructure:"filename"`
	Level      string `mapstructure:"level"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type MysqlConfig struct {
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
	User    string `mapstructure:"user"`
	Pass    string `mapstructure:"pass"`
	DB      string `mapstructure:"dbname"`
	MaxOpen int    `mapstructure:"max_open"`
	MaxIdle int    `mapstructure:"max_idle"`
}
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Pass     string `mapstructure:"pass"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

func Init(filePath string) (err error) {
	// 方式1：指定配置文件
	viper.SetConfigFile(filePath) // 指定配置文件

	// 方式2：
	//viper.SetConfigName("conf") // 配置文件名称(无扩展名)
	//viper.SetConfigType("yaml") // 指定文件类型（专用来从远程获取配置信息时指定配置文件类型，本地不生效，配合配置中心第三方使用）
	//viper.AddConfigPath(".")   // 查找配置文件所在的路径
	err = viper.ReadInConfig() // 查找并读取配置文件
	if err != nil {            // 处理读取配置文件的错误
		fmt.Println("viper.ReadInConfig failed,err")
		return err
	}

	// 把读取到的配置信息反序列化到conf变量中
	if err = viper.Unmarshal(Conf); err != nil {
		fmt.Println("unmarshal conf failed")
		return err
	}

	// 热加载配置
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置文件被修改啦...")
		if err = viper.Unmarshal(Conf); err != nil {
			fmt.Printf("unmarshal conf failed")
		}
	})

	return
}
