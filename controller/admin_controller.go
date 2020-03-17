package controller

import (
	"CmsProject/service"
	"CmsProject/utils"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
	"strconv"
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
	ADMIN = "admin"
)

func (ac *AdminController) GetAll() mvc.Result {
	limit, err := strconv.Atoi(ac.Ctx.FormValue("limit"))
	offset, err := strconv.Atoi(ac.Ctx.FormValue("offset"))
	if err != nil {
		offset = 0
		limit = 20
	}
	admin, err := ac.Service.GetAdminAll(limit, offset)
	if err != nil {
		panic(err.Error())
		return nil
	}
	return mvc.Response{
		Object: admin,
	}
}

/**
 * 管理员退出功能
 * 请求类型：Get
 * 请求url：admin/singout
 */
func (ac *AdminController) GetSingout() mvc.Result {

	//删除session，下次需要从新登录
	ac.Session.Delete(ADMIN)
	return mvc.Response{
		Object: map[string]interface{}{
			"status":  utils.RECODE_OK,
			"success": utils.Recode2Text(utils.RESPMSG_SIGNOUT),
		},
	}
}

/**
 * 处理获取管理员总数的路由请求
 * 请求类型：Get
 * 请求Url：admin/count
 */
func (ac *AdminController) GetCount() mvc.Result {

	count, err := ac.Service.GetAdminCount()
	if err != nil {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  utils.RECODE_FAIL,
				"message": utils.Recode2Text(utils.RESPMSG_ERRORADMINCOUNT),
				"count":   0,
			},
		}
	}
	return mvc.Response{
		Object: map[string]interface{}{
			"status": utils.RECODE_OK,
			"count":  count,
		},
	}
}

/**
 * 获取管理员信息接口
 * 请求类型：Get
 * 请求url：/admin/info
 */
func (ac *AdminController) GetInfo() mvc.Result {

	//从session中获取信息
	userByte := ac.Session.Get(ADMIN)

	//session为空
	if userByte == nil {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  utils.RECODE_UNLOGIN,
				"type":    utils.EEROR_UNLOGIN,
				"message": utils.Recode2Text(utils.EEROR_UNLOGIN),
			},
		}
	}

	adminId, err := ac.Session.GetInt64(ADMIN)

	//解析失败
	if err != nil {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  utils.RECODE_UNLOGIN,
				"type":    utils.EEROR_UNLOGIN,
				"message": utils.Recode2Text(utils.EEROR_UNLOGIN),
			},
		}
	}
	adminObject, exit := ac.Service.GetByAdminId(adminId)

	if !exit {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  "0",
				"success": "登录失败",
				"message": "用户名或者密码错误,请重新登录",
			},
		}
	}

	//解析成功
	return mvc.Response{
		Object: map[string]interface{}{
			"status": utils.RECODE_OK,
			"data":   adminObject.AdminToRespDesc(),
		},
	}
}

/**
 * 管理员登录功能
 * 接口：/admin/login
 */
func (ac *AdminController) PostLogin() mvc.Result {

	var adminLogin AdminLogin
	ac.Ctx.ReadJSON(&adminLogin)

	//数据参数检验
	if adminLogin.UserName == "" || adminLogin.PassWord == "" {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  "0",
				"success": "登录失败",
				"message": "用户名或密码为空,请重新填写后尝试登录",
			},
		}
	}

	//根据用户名、密码到数据库中查询对应的管理信息
	admin, exist := ac.Service.GetByAdminNameAndPassword(adminLogin.UserName, adminLogin.PassWord)

	//管理员不存在
	if !exist {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  "0",
				"success": "登录失败",
				"message": "用户名或者密码错误,请重新登录",
			},
		}
	}

	//管理员存在 设置session
	ac.Session.Set(ADMIN, admin.AdminId)

	//管理员存在 设置session
	return mvc.Response{
		Object: map[string]interface{}{
			"status":  "1",
			"success": "登录成功",
			"message": "管理员登录成功",
		},
	}
}
