package controller

import (
	"CmsProject/service"
	"CmsProject/utils"
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
	"strconv"
)

/**
 * 用户控制器结构体：用来实现处理用户模块的接口的请求，并返回给客户端
 */
type UserController struct {
	Session *sessions.Session
	Service service.UserService

	//上下文对象
	Ctx iris.Context
}

//每一页最大的内容
const MaxLimit = 50

/**
 * 获取所有用户的信息
 * 请求类型：Get
 * 请求Url：/v1/users/list
 */
func (uc *UserController) GetList() mvc.Result {
	offsetStr := uc.Ctx.FormValue("offset")
	limitStr := uc.Ctx.FormValue("limit")
	var offset int
	var limit int

	//判断offset和limit两个变量任意一个都不能为""
	if offsetStr == "" || limitStr == "" {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_ERROR_USERLIST,
				"message": utils.Recode2Text(utils.RESPMSG_ERROR_USERLIST),
			},
		}
	}

	offset, err := strconv.Atoi(offsetStr)
	limit, err = strconv.Atoi(limitStr)
	if err != nil {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_ERROR_USERLIST,
				"message": utils.Recode2Text(utils.RESPMSG_ERROR_USERLIST),
			},
		}
	}

	//做页数的限制检查
	if offset <= 0 {
		offset = 0
	}

	//做最大的限制
	if limit > MaxLimit {
		limit = MaxLimit
	}
	userList := uc.Service.GetUserList(offset, limit)
	if len(userList) == 0 {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_ERROR_USERLIST,
				"message": utils.Recode2Text(utils.RESPMSG_ERROR_USERLIST),
			},
		}
	}

	//将查询到的用户数据进行转换成前端需要的内容
	var respList []interface{}
	for _, user := range userList {
		respList = append(respList, user.UserToRespDesc())
	}

	//返回用户列表
	return mvc.Response{
		Object: &respList,
	}
}

/**
 * 获取用户总数
 * 请求类型：Get
 * 请求Url：/v1/users/count
 */
func (uc *UserController) GetCount() mvc.Result {
	//用户总数
	total, err := uc.Service.GetUserTotalCount()

	//请求出现错误
	if err != nil {
		return mvc.Response{
			Object: map[string]interface{}{
				"status": utils.RECODE_FAIL,
				"count":  0,
			},
		}
	}
	fmt.Println("注册用户总数:", total)

	//正常情况的返回值
	return mvc.Response{
		Object: map[string]interface{}{
			"status": utils.RECODE_OK,
			"count":  total,
		},
	}
}
