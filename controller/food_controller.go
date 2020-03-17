package controller

import (
	"CmsProject/service"
	"CmsProject/utils"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"strconv"
)

type FoodController struct {
	Ctx     iris.Context
	Service service.FoodService
}

func (fc *FoodController) Get() mvc.Result {
	offset, err := strconv.Atoi(fc.Ctx.Params().Get("offset"))
	limit, err := strconv.Atoi(fc.Ctx.Params().Get("limit"))
	if err != nil {
		offset = 0
		limit = 20
	}
	list, err := fc.Service.GetFoodList(offset, limit)
	if err != nil {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_ERROR_FOODLIST,
				"message": utils.Recode2Text(utils.RESPMSG_ERROR_FOODLIST),
			},
		}
	}
	//成功
	return mvc.Response{
		Object: &list,
	}
}

func (fc *FoodController) GetCount() mvc.Result {
	result, err := fc.Service.GetFoodCount()
	if err != nil {
		return mvc.Response{
			Object: map[string]interface{}{
				"status": utils.RESPMSG_FAIL,
				"count":  0,
			},
		}
	}

	return mvc.Response{
		Object: map[string]interface{}{
			"status": utils.RESPMSG_OK,
			"count":  result,
		},
	}
}
