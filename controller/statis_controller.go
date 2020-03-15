package controller

import (
	"CmsProject/service"
	"CmsProject/utils"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
	"strings"
	"time"
)

/**
 * 统计功能控制者
 */
type StatisController struct {
	//上下文环境对象
	Ctx iris.Context

	Session *sessions.Session

	//统计功能的服务实现接口
	Service service.StatisService
}

const (
	ADMINMODULE = "ADMIN_"
	USERMODULE  = "USER_"
	ORDERMODULE = "ORDER_"
)

/**
 * 解析统计功能路由请求
 */
func (sc *StatisController) GetCount() mvc.Result {
	path := sc.Ctx.Path()
	var pathSlice []string
	if path != "" {
		pathSlice = strings.Split(path, "/")
	}

	//不符合请求格式
	if len(pathSlice) != 5 {
		return mvc.Response{
			Object: map[string]interface{}{
				"status": utils.RECODE_FAIL,
				"count":  0,
			},
		}
	}

	//将最前面的"/"去掉
	pathSlice = pathSlice[1:]
	model := pathSlice[1]
	date := pathSlice[2]
	var result int64
	switch model {
	case "user":
		if date == "NaN-NaN-NaN" {
			date = time.Now().Format("2006-01-02")
		}
		userResult := sc.Session.Get(USERMODULE + date)
		if userResult != nil {
			userResult = userResult.(float64)
			return mvc.Response{
				Object: map[string]interface{}{
					"status": utils.RECODE_OK,
					"count":  userResult,
				},
			}
		} else {
			iris.New().Logger().Info(date)
			result = sc.Service.GetUserDailyCount(date)

			//如果不是当日时间，就设置缓存
			if date != time.Now().Format("2006-01-02") {
				sc.Session.Set(USERMODULE+date, result)
			}
		}
	case "order":
		if date == "NaN-NaN-NaN" {
			date = time.Now().Format("2006-01-02")
		}
		orderResult := sc.Session.Get(ORDERMODULE + date)
		if orderResult != nil {
			orderResult = orderResult.(float64)
			return mvc.Response{
				Object: map[string]interface{}{
					"status": utils.RECODE_OK,
					"count":  orderResult,
				},
			}
		} else {
			iris.New().Logger().Info(date)
			result = sc.Service.GetOrderDailyCount(date)

			//设置缓存
			if date != time.Now().Format("2006-01-02") {
				sc.Session.Set(ORDERMODULE+date, result)
			}
		}
	case "admin":
		if date == "NaN-NaN-NaN" {
			date = time.Now().Format("2006-01-02")
		}
		adminResult := sc.Session.Get(ADMINMODULE + date)
		if adminResult != nil {
			adminResult = adminResult.(float64)
			return mvc.Response{
				Object: map[string]interface{}{
					"status": utils.RECODE_OK,
					"count":  adminResult,
				},
			}
		} else {
			iris.New().Logger().Info(date)
			result = sc.Service.GetAdminDailyCount(date)

			//设置缓存
			if date != time.Now().Format("2006-01-02") {
				sc.Session.Set(ADMINMODULE+date, result)
			}
		}
	}
	return mvc.Response{
		Object: map[string]interface{}{
			"status": utils.RECODE_OK,
			"count":  result,
		},
	}
}
