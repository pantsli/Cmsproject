package controller

import (
	"CmsProject/service"
	"encoding/json"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
)

/**
 * 管理员控制器
 */
type AdminController struct {
	//iris框架自动为每个请求都绑定上下文对象
	Ctx iris.Context

	//session对象
	Session *sessions.Session

	//admin功能实体
	Service service.AdminService
}

type AdminLogin struct {
	UserName string `json:"user_name"`
	PassWord string `json:"password"`
}

const (
	ADMIN="admin"
)

/**
 * 管理员登录功能
 * 接口：/admin/login
 */
func (ac *AdminController)PostLogin()mvc.Result{

	var adminLogin AdminLogin
	ac.Ctx.ReadJSON(&adminLogin)

	//数据参数检验
	if adminLogin.UserName==""||adminLogin.PassWord==""{
		return mvc.Response{
			Object:map[string]interface{}{
				"status":  "0",
				"success": "登录失败",
				"message": "用户名或密码为空,请重新填写后尝试登录",
			},
		}
	}

	//根据用户名、密码到数据库中查询对应的管理信息
	admin,exist:=ac.Service.GetByAdminNameAndPassword(adminLogin.UserName,adminLogin.PassWord)

	//管理员不存在
	if !exist{
		return mvc.Response{
			Object:map[string]interface{}{
				"status":  "0",
				"success": "登录失败",
				"message": "用户名或者密码错误,请重新登录",
			},
		}
	}

	//管理员存在 设置session
	userByte,_:=json.Marshal(admin)
	ac.Session.Set(ADMIN,userByte)

	//管理员存在 设置session
	return mvc.Response{
		Object: map[string]interface{}{
			"status":  "1",
			"success": "登录成功",
			"message": "管理员登录成功",
		},
	}
}