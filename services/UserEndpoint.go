package services

import (
	"context"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/go-kit/kit/endpoint"
)

type UserRequest struct {
	Uid    int64  `json:"uid"`
	Method string `json:"method"`
}

type UserResponse struct {
	Result string `json:"result"`
}

//MakeUserEndpoint 创建endpoint
func MakeUserEndpoint(service IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*UserRequest)
		result := ""
		if !ok {
			logs.Error("the user is not UserRequest")
			result = "the request is not UserRequest"
		}
		if service == nil {
			logs.Error("the services is nil")
			result = "the service is nil"
		}
		if req.Method == "GET" {
			logs.Debug("the method is get")
			result = service.GetUserName(req.Uid)
		} else if req.Method == "DELETE" {
			logs.Debug("the method is delete")
			err = service.DeleteUser(req.Uid)
			if err != nil {
				result = err.Error()
			} else {
				result = fmt.Sprintf("delete user:%d success", req.Uid)
			}
		}
		return &UserResponse{Result: result}, nil
	}
}
