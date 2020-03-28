package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/astaxie/beego/logs"
	http2 "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/jaydenwen123/gokit-practice-server/services"
)

var port int
var name string
var id string

func main() {

	//从控制台接收参数
	flag.IntVar(&port, "port", 3456, "please the server listen port")
	flag.StringVar(&name, "name", "", "please the server register name")
	flag.Parse()
	if port == 0 {
		panic("the port must not empty...")
	}

	if name == "" || len(name) == 0 {
		panic("the name must not empty...")
	}

	id = genServiceId(name)

	userService := &services.UserService{}
	endpoint := services.MakeUserEndpoint(userService)
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
	errChan:=make(chan  error,1)
	go func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM,syscall.SIGKILL)
		errChan<-fmt.Errorf("%v",<-signalChan)
	}()

	go func() {
		//注册服务
		//RegisterService("myserver2","userService","127.0.0.1",2346)
		RegisterService(id, name, "127.0.0.1", port)
		//剔除服务
		logs.Debug("the server is starting.....")
		err := http.ListenAndServe(fmt.Sprintf(":%d", port), router)
		if err != nil {
			logs.Error("the server started error:", err.Error())
			errChan<-err
		}
	}()

	<-errChan
	logs.Error("the server is exited....")
	DeRegisterService(id)
	logs.Debug("DeRegisterService from the consul ....")
}
