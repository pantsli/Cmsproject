package controller

import (
	"CmsProject/service"
	"CmsProject/utils"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
	"strconv"
)

//商店功能模块控制结构体
type ShopController struct {
	Session *sessions.Session

	//上下文对象
	Ctx     iris.Context
	Service service.ShopService
}

/**
 * 获取商铺列表
 * 请求Url：/shopping/restaurants
 * 请求类型：get
 */
func (sc *ShopController) Get() mvc.Result {
	offsetStr := sc.Ctx.FormValue("offset")
	limitStr := sc.Ctx.FormValue("limit")

	if offsetStr == "" || limitStr == "" { //设置默认值
		offsetStr = "0"
		limitStr = "20"
	}
	offset, err := strconv.Atoi(offsetStr)
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		offset = 0
		limit = 20
	}
	shopList := sc.Service.GetShopList(offset, limit)

	if len(shopList) < 0 {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_ERROR_RESTLIST,
				"message": utils.Recode2Text(utils.RESPMSG_ERROR_RESTLIST),
			},
		}
	}

	var respList []interface{}

	for _, shop := range shopList {
		respList = append(respList, shop.ShopToRespDesc())
	}
	return mvc.Response{
		Object: respList,
	}
}
func (sc *ShopController) GetCount() mvc.Result {
	result, err := sc.Service.GetShopCount()
	if err != nil {
		return mvc.Response{
			Object: map[string]interface{}{
				"status": utils.RECODE_FAIL,
				"count":  0,
			},
		}
	}

	//正常情况的返回值
	return mvc.Response{
		Object: map[string]interface{}{
			"status": utils.RECODE_OK,
			"count":  result,
		},
	}
}
