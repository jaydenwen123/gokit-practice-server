package services

import "fmt"

type IUserService interface {
	GetUserName(uid int64) string
	DeleteUser(uid int64) error
}

type UserService struct {
}

//GetUserName 获取用户名称
func (u *UserService) GetUserName(uid int64) string {
	if uid == 111 {
		return "wenxiaofei"
	}
	return "what funny"
}

//DeleteUser 删除用户
func (u *UserService) DeleteUser(uid int64) error {
	if uid == 111 {
		return fmt.Errorf("there is not permitted to delete administator")
	}
	return nil
}
