package service

import (
	"CmsProject/model"
	"github.com/go-xorm/xorm"
)

/**
 * 用户模块功能服务接口
 */
type UserService interface {
	//获取用户总数
	GetUserTotalCount() (int64, error)
	//用户列表
	GetUserList(offset, limit int) []*model.User
}

/**
 * 用户服务实现结构体
 */
type userService struct {
	engine *xorm.Engine
}

/**
 * 实例化用户服务结构实体对象
 */
func NewUserService(db *xorm.Engine) UserService {
	return &userService{engine: db}
}

/**
 * 请求用户列表数据
 * offset：偏移数量
 * limit：一次请求获取的数据条数
 */
func (us *userService) GetUserList(offset, limit int) []*model.User {
	var userList []*model.User

	err := us.engine.Where("del_flag=?", 0).Limit(limit, offset).Find(&userList)

	if err != nil {
		panic(err.Error())
		return nil
	}

	return userList
}

/**
 * 请求总的用户数量
 * 返回值：总用户数量
 */
func (us *userService) GetUserTotalCount() (int64, error) {

	//查询del_flag 为0 的用户的总数量；del_flag:0 正常状态；del_flag:1 用户注销或者被删除
	count, err := us.engine.Where("del_flag = ?", 0).Count(new(model.User))
	if err != nil {
		panic(err.Error())
		return 0, err
	}

	//用户总数
	return count, nil
}
