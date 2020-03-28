package main

import (
	"fmt"
	"github.com/jaydenwen123/go-util"
	"net/http"

	"github.com/astaxie/beego/logs"
	http2 "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/jaydenwen123/gokit-practice-server/services"

	consulapi "github.com/hashicorp/consul/api"
)

func main() {
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
	logs.Debug("the server is starting.....")
	RegisterService("myserver2","userService","127.0.0.1",2346)
	http.ListenAndServe(":2346", router)
}

//服务注册
func RegisterService( id, name,address string, port int,) {
	client := makeClient()
	agent := client.Agent()
	reg := &consulapi.AgentServiceRegistration{
		Name:    name,
		ID:      id,
		Tags:    []string{"primary", "backup"},
		Address: address,
		Port:    2346,
		Check: &consulapi.AgentServiceCheck{
			HTTP: fmt.Sprintf("http://%s:%d/health",address,port),
			Interval:"10s",
		},
	}
	if err := agent.ServiceRegister(reg); err != nil {
		logs.Error("ServiceRegister err: %v", err)
	}

	services, err := agent.Services()
	if err != nil {
		logs.Error("err: %v", err)
	}
	logs.Debug("the services:%s", util.Obj2JsonStr(services))
	//if _, ok := services["myserver2"]; !ok {
	//logs.Error("missing service: %#v", services)
	//}
	checks, err := agent.Checks()

	if err != nil {
		logs.Error("get checks err: %v", err)
	}
	logs.Debug("the checks:%s", util.Obj2JsonStr(checks))
	//chk, ok := checks["service:foo"]
	//if !ok {
	//	logs.Error("missing check: %v", checks)
	//}

	// Checks should default to critical
	//if chk.Status != consulapi.HealthCritical {
	//	logs.Error("Bad: %#v", chk)
	//}

	//state, out, err := agent.AgentHealthServiceByID("foo2")
	//require.Nil(t, err)
	//require.Nil(t, out)
	//require.Equal(t, HealthCritical, state)
	//
	//state, out, err = agent.AgentHealthServiceByID("foo")
	//require.Nil(t, err)
	//require.NotNil(t, out)
	//require.Equal(t, HealthCritical, state)
	//require.Equal(t, 8000, out.Service.Port)
	//
	//state, outs, err := agent.AgentHealthServiceByName("foo")
	//require.Nil(t, err)
	//require.NotNil(t, outs)
	//require.Equal(t, HealthCritical, state)
	//require.Equal(t, 8000, outs[0].Service.Port)

	if err = agent.ServiceDeregister("foo"); err != nil {
		logs.Error("the ServiceDeregister failed  err: %v", err)
	} else {
		logs.Debug("the ServiceDeregister success...")
	}
}

func makeClient() *consulapi.Client {

	return makeClientWithConfig(nil)
}

type configCallback func(c *consulapi.Config)

func makeClientWithConfig(cb1 configCallback) *consulapi.Client {
	// Make client config
	conf := consulapi.DefaultConfig()
	if cb1 != nil {
		cb1(conf)
	}

	// Create client
	client, err := consulapi.NewClient(conf)
	if err != nil {
		logs.Error("new consul Client  err: %v", err)
	}
	return client
}
