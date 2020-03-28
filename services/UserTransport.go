package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

//DecodeRequest 解码
func DecodeRequest(ctx context.Context, req *http.Request) (request interface{}, err error) {
	//从request中获取参数
	//改造这块
	//uidStr := req.URL.Query().Get("uid")
	logs.Debug("the port:",req.URL.Port())
	vars := mux.Vars(req)
	uidStr, ok := vars["uid"]
	uid := int64(0)
	if ok && uidStr != "" {
		uid, err = strconv.ParseInt(uidStr, 10, 64)
		if err != nil {
			logs.Error("the uid is not a number.")
			return nil, err
		}
		return &UserRequest{Uid: uid, Method: req.Method}, nil
	}
	return nil, fmt.Errorf("the uid is not valid")
}

//EncodeResponsefunc 编码
func EncodeResponsefunc(ctx context.Context, writer http.ResponseWriter, response interface{}) error {
	writer.Header().Set("Content-Type", "Application/json")
	encoder := json.NewEncoder(writer)
	if response == nil {
		return encoder.Encode("the op is failed")
	}
	return encoder.Encode(response)
}
