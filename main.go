package main

import (
	"flag"
	"fmt"
	"net/http"

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
	flag.IntVar(&port,"port",3456,"please the server listen port")
	flag.StringVar(&name,"name","","please the server register name")
	flag.Parse()
	if port==0 {
		panic("the port must not empty...")
	}

	if name == "" || len(name)==0	 {
		panic("the name must not empty...")
	}

	id=genServiceId(name)

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
	//注册服务
	//RegisterService("myserver2","userService","127.0.0.1",2346)
	RegisterService(id,name,"127.0.0.1",port)

	logs.Debug("the server is starting.....")
	http.ListenAndServe(fmt.Sprintf(":%d",port), router)
}


