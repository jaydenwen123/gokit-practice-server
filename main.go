package main

import (
	"net/http"

	"github.com/astaxie/beego/logs"
	http2 "github.com/go-kit/kit/transport/http"
	"github.com/jaydenwen123/gokit-practice-server/services"
	"github.com/gorilla/mux"
)

func main() {
	userService := &services.UserService{}
	endpoint := services.MakeUserEndpoint(userService)
	handler := http2.NewServer(endpoint, services.DecodeRequest, services.EncodeResponsefunc)
	router := mux.NewRouter()
	router.Handle("/user/{uid:\\d+}",handler)

	//添加心跳检测
	router.HandleFunc("/health", func(rspWriter http.ResponseWriter, request *http.Request) {
		logs.Debug("there is the heart check ....")
		rspWriter.Header().Set("Content-type","application-json")
		rspWriter.Write([]byte(`{"status":"ok"}`))
	})

	//router.Methods("GET","DELETE")
	logs.Debug("the server is starting.....")
	http.ListenAndServe(":2345", router)
}
