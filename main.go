package main

import (
	"flag"
	"fmt"
	"github.com/jaydenwen123/gokit-practice-server/consul_util"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/astaxie/beego/logs"
	http2 "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/jaydenwen123/gokit-practice-server/services"
)

var Port int
var Name string
var Id string

func main() {

	//从控制台接收参数
	flag.IntVar(&Port, "port", 3456, "please the server listen Port")
	flag.StringVar(&Name, "name", "", "please the server register Name")
	flag.Parse()
	if Port == 0 {
		panic("the Port must not empty...")
	}

	if Name == "" || len(Name) == 0 {
		panic("the Name must not empty...")
	}

	Id = consul_util.GenServiceId(Name)

	userService := &services.UserService{}
	endpoint := services.MakeUserEndpoint(userService,Port)
	handler := http2.NewServer(endpoint, services.DecodeRequest, services.EncodeResponsefunc)
	router := mux.NewRouter()
	router.Handle("/user/{uid:\\d+}", handler)

	//添加心跳检测
	router.HandleFunc("/health", func(rspWriter http.ResponseWriter, request *http.Request) {
		logs.Debug("there is the heart check ....")
		rspWriter.Header().Set("Content-type", "application-json")
		rspWriter.Write([]byte(`{"status":"ok"}`))
	})

	//router.Methods("GET","DELETE")

	//服务退出有两种情况
	//1和2通过监听信号实现
	//1.CTRL+C终止程序
	//2.kill -9 +pid
	//直接剔除服务
	//3.程序启动错误
	errChan := make(chan error, 1)
	go func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
		errChan <- fmt.Errorf("%v", <-signalChan)
	}()

	go func() {
		//注册服务
		//RegisterService("myserver2","userService","127.0.0.1",2346)
		consul_util.RegisterService(Id, Name, "127.0.0.1", Port)
		//剔除服务
		logs.Debug("the server is starting.....")
		err := http.ListenAndServe(fmt.Sprintf(":%d", Port), router)
		if err != nil {
			logs.Error("the server started error:", err.Error())
			errChan <- err
		}
	}()

	<-errChan
	logs.Error("the server is exited....")
	consul_util.DeRegisterService(Id)
	logs.Debug("DeRegisterService from the consul ....")
}
