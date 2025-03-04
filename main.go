package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"web_app/controller"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/logger"
	"web_app/pkg/snowflake"
	"web_app/routes"
	"web_app/settings"
)

// @title 这里写标题
// @version 1.0
// @description 这里写描述信息
// @termsOfService http://swagger.io/terms/

// @contact.name 这里写联系人信息
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 127.0.0.1:8081
// @BasePath /api/v1
func main() {
	// 定义命令行参数
	var filePath string
	flag.StringVar(&filePath, "f", "conf.yaml", "配置文件路径")
	// 解析命令行参数
	flag.Parse()

	// 1.加载配置
	if err := settings.Init(filePath); err != nil {
		fmt.Println("Init settings failed，err:", err)
		return
	}

	// 2.初始化日志
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Println("Init logger failed，err:", err)
		return
	}
	defer zap.L().Sync()
	zap.L().Debug("日志初始化成功...")

	// 3.初始化Mysql连接
	if err := mysql.Init(settings.Conf.MysqlConfig); err != nil {
		fmt.Println("Init mysql failed，err:", err)
		return
	}
	defer mysql.Close()
	zap.L().Debug("Mysql初始化成功...")

	// 4.初始化Redis连接
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Println("Init redis failed，err:", err)
		return
	}
	defer redis.Close()
	zap.L().Debug("Redis初始化成功...")

	// 雪花算法生成唯一用户ID初始化
	err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineID)
	if err != nil {
		fmt.Println("Init snowflake failed，err:", err)
		return
	}

	// 初始化gin框架内置的校验翻译器
	err = controller.InitTrans("zh")
	if err != nil {
		return
	}

	// 5.注册路由
	r := routes.SetUp()

	// 6.启动服务（优雅关机）
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Conf.AppConfig.Port),
		Handler: r,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}

	zap.L().Info("Server exiting")

}
