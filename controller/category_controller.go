package controller

import (
	"CmsProject/model"
	"CmsProject/service"
	"CmsProject/utils"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
	"strconv"
)

/**
 * 食品种类控制器
 */
type CategoryController struct {
	Session *sessions.Session
	Ctx     iris.Context
	Service service.CategoryService
}

/**
 * 添加食品种类实体
 */
type CategoryEntity struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	RestaurantId string `json:"restaurant_id"`
}

func (cc *CategoryController) BeforeActivation(a mvc.BeforeActivation) {
	a.Handle("POST", "/addcategory", "AddCategory")
	a.Handle("GET", "/getcategory/{shopId}", "GetCategoryByShopId")
}

/**
 * url：/shopping/getcategory/1
 * type：get
 * desc：根据商铺Id获取对应的商铺的食品种类列表信息
 */
func (cc *CategoryController) GetCategoryByShopId() mvc.Result {
	shopIdStr := cc.Ctx.Params().Get("shopId")
	if shopIdStr == "" {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  utils.RECODE_FAIL,
				"message": utils.Recode2Text(utils.RESPMSG_ERROR_CATEGORIES),
			},
		}
	}
	shopId, err := strconv.Atoi(shopIdStr)
	if err != nil {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  utils.RECODE_FAIL,
				"message": utils.Recode2Text(utils.RESPMSG_ERROR_CATEGORIES),
			},
		}
	}

	//调用服务实体功能类查询商铺对应的食品种类列表
	categories, err := cc.Service.GetCategoryByShopId(int64(shopId))
	if err != nil {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  utils.RECODE_FAIL,
				"message": utils.Recode2Text(utils.RESPMSG_ERROR_CATEGORIES),
			},
		}
	}

	//返回对应店铺的食品种类类型
	return mvc.Response{
		Object: map[string]interface{}{
			"status":        utils.RECODE_OK,
			"category_list": &categories,
		},
	}
}

/**
 * url：/shopping/addcategory
 * type：post
 * desc：添加食品种类记录
 */
func (cc *CategoryController) AddCategory() mvc.Result {
	var categoryEntity *CategoryEntity
	cc.Ctx.ReadJSON(&categoryEntity)
	if categoryEntity.Name == "" {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  utils.RECODE_FAIL,
				"message": utils.Recode2Text(utils.RESPMSG_ERROR_CATEGORYADD),
			},
		}
	}

	restaurantId, _ := strconv.Atoi(categoryEntity.RestaurantId)

	//构造要添加的数据记录
	foodCategory := &model.FoodCategory{
		CategoryName:     categoryEntity.Name,
		CategoryDesc:     categoryEntity.Description,
		Level:            1,
		ParentCategoryId: 0,
		RestaurantId:     int64(restaurantId),
	}
	bool := cc.Service.AddCategory(foodCategory)
	if !bool {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  utils.RECODE_FAIL,
				"message": utils.Recode2Text(utils.RESPMSG_ERROR_CATEGORYADD),
			},
		}
	}

	//成功
	return mvc.Response{
		Object: map[string]interface{}{
			"status":  utils.RECODE_OK,
			"message": utils.Recode2Text(utils.RESPMSG_SUCCESS_CATEGORYADD),
		},
	}
}
