package main

import "fmt"
import "os"
import "os/signal"
import "syscall"
import "logger"
import "flag"
import "runtime"
import "server"


var configFile *string = flag.String("config", "/etc/mag.yaml", "server config file")
var	srv *server.Server
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	//读取配置文件信息
	flag.Parse()

	sc:=make(chan os.Signal,1)
	signal.Notify(sc,syscall.SIGTERM,syscall.SIGINT,syscall.SIGQUIT,syscall.SIGPIPE)

	go func() {
		for{
			sig:=<-sc
			if sig==syscall.SIGQUIT ||sig==syscall.SIGTERM {
				logger.Info("main","main","server will stop service! ",0,"revice signal:",sig)
				Stop()
			}

			if sig==syscall.SIGINT{
				Start()
			}
		}

	}()


	fmt.Println("hello world!")
}

func LoadConfig() {
	logger.Info("main","main","server will stop service! ",0,"load configure file!")

}
func InitLog() {
	logger.Info("main","main","server will stop service! ",0,"init logger module!")
}

func Start()  {
	logger.Info("main","main","server will stop service! ",0,"start service!")

	//加载配置文件
	LoadConfig()

	//初始化Logger
	InitLog()

	srv.StartServer()
}

func Stop()  {
	logger.Info("main","main","server will stop service! ",0,"stop service!")

}