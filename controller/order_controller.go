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

type OrderController struct {
	Session *sessions.Session
	Service service.OrderService
	Ctx     iris.Context
}

/**
 * 获取订单列表
 */
func (oc *OrderController) Get() mvc.Result {
	offsetStr := oc.Ctx.FormValue("offset")
	limitStr := oc.Ctx.FormValue("limit")

	var offset int
	var limit int

	//判断offset和limit两个变量任意一个都不能为""
	if offsetStr == "" || limitStr == "" {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_ERROR_ORDERLIST,
				"message": utils.Recode2Text(utils.RESPMSG_ERROR_ORDERLIST),
			},
		}
	}

	offset, err := strconv.Atoi(offsetStr)
	limit, err = strconv.Atoi(limitStr)

	if err != nil {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_ERROR_ORDERLIST,
				"message": utils.Recode2Text(utils.RESPMSG_ERROR_ORDERLIST),
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

	orderList := oc.Service.GetOrderList(offset, limit)
	if len(orderList) == 0 {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_ERROR_ORDERLIST,
				"message": utils.Recode2Text(utils.RESPMSG_ERROR_ORDERLIST),
			},
		}
	}

	//将查询到的用户数据进行转换成前端需要的内容
	var respList []interface{}
	for _, order := range orderList {
		respList = append(respList, order.OrderDetail2Resp())
	}

	//返回用户列表
	return mvc.Response{
		Object: &respList,
	}
}

/**
 * 查询订单记录总数
 */
func (oc *OrderController) GetCount() mvc.Result {
	count, err := oc.Service.GetCount()
	if err != nil {
		return mvc.Response{
			Object: map[string]interface{}{
				"status": utils.RECODE_FAIL,
				"count":  0,
			},
		}
	}
	fmt.Println("订单总数:", count)
	return mvc.Response{
		Object: map[string]interface{}{
			"status": utils.RECODE_OK,
			"count":  count,
		},
	}
}
