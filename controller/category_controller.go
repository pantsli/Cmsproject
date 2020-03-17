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

/**
 * 即将添加的食品记录实体
 */
type AddFoodEntity struct {
	Name         string   `json:"name"`          //食品名称
	Description  string   `json:"description"`   //食品描述
	ImagePath    string   `json:"image_path"`    //食品图片地址
	Activity     string   `json:"activity"`      //食品活动
	Attributes   []string `json:"attributes"`    //食品特点
	Specs        []Specs  `json:"specs"`         //食品规格
	CategoryId   int      `json:"category_id"`   //食品种类  种类id
	RestaurantId string   `json:"restaurant_id"` //哪个店铺的食品 店铺id
}

//食品规格
type Specs struct {
	Specs      string `json:"specs"`
	PackingFee int    `json:"packing_fee"`
	Price      int    `json:"price"`
}

func (cc *CategoryController) BeforeActivation(a mvc.BeforeActivation) {
	a.Handle("POST", "/addcategory", "AddCategory")
	a.Handle("GET", "/getcategory/{shopId}", "GetCategoryByShopId")
	a.Handle("POST", "/addfood", "PostAddFood")
	a.Handle("DELETE", "/v2/food/{food_id}", "Delfood")
	a.Handle("POST", "/addShop", "PostAddShop")
	a.Handle("DELETE", "/restaurant/{shop_id}", "DelShop")
}

/**
 * 删除商户记录
 *
 */
func (cc *CategoryController) DelShop() mvc.Result {
	shopId, err := strconv.Atoi(cc.Ctx.Params().Get("shop_id"))
	if err != nil {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_HASNOACCESS,
				"message": utils.Recode2Text(utils.RESPMSG_HASNOACCESS),
			},
		}
	}
	isBool := cc.Service.DeleteShop(shopId)
	if !isBool {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_HASNOACCESS,
				"message": utils.Recode2Text(utils.RESPMSG_HASNOACCESS),
			},
		}
	}
	return mvc.Response{
		Object: map[string]interface{}{
			"status":  utils.RECODE_OK,
			"type":    utils.RESPMSG_SUCCESS_DELETESHOP,
			"message": utils.Recode2Text(utils.RESPMSG_SUCCESS_DELETESHOP),
		},
	}
}

func (cc *CategoryController) PostAddShop() mvc.Result {
	var shop model.Shop
	err := cc.Ctx.ReadJSON(&shop)
	if err != nil {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  utils.RECODE_FAIL,
				"message": utils.Recode2Text(utils.RESPMSG_FAIL_ADDREST),
			},
		}
	}

	saveShop := cc.Service.SaveShop(shop)
	if !saveShop {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  utils.RECODE_FAIL,
				"message": utils.Recode2Text(utils.RESPMSG_FAIL_ADDREST),
			},
		}
	}
	return mvc.Response{
		Object: map[string]interface{}{
			"status":     utils.RECODE_OK,
			"message":    utils.Recode2Text(utils.RESPMSG_SUCCESS_ADDREST),
			"shopDetail": shop,
		},
	}

}

func (cc *CategoryController) Delfood() mvc.Result {
	foodId, err := strconv.Atoi(cc.Ctx.Params().Get("food_id"))
	if err != nil {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_ERROR_FOODDELE,
				"message": utils.Recode2Text(utils.RESPMSG_ERROR_FOODDELE),
			},
		}
	}

	isSuccess := cc.Service.DelFood(foodId)
	if !isSuccess {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_ERROR_FOODDELE,
				"message": utils.Recode2Text(utils.RESPMSG_ERROR_FOODDELE),
			},
		}
	}

	return mvc.Response{
		Object: map[string]interface{}{
			"status":  utils.RECODE_OK,
			"success": utils.Recode2Text(utils.RESPMSG_SUCCESS_FOODDELE),
		},
	}
}

/**
 * url: /shopping/addfood
 * type：post
 * descs：添加食品数据功能
 */
func (cc *CategoryController) PostAddFood() mvc.Result {
	var foodEntity AddFoodEntity
	err := cc.Ctx.ReadJSON(&foodEntity)
	if err != nil {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_ERROR_FOODADD,
				"message": utils.Recode2Text(utils.RESPMSG_ERROR_FOODADD),
			},
		}
	}
	var food model.Food
	food.Name = foodEntity.Name
	food.Name = foodEntity.Name
	food.Description = foodEntity.Description
	food.ImagePath = foodEntity.ImagePath
	food.Activity = foodEntity.Activity
	food.CategoryId = int64(foodEntity.CategoryId)
	var Attributes string
	for _, food := range foodEntity.Attributes {
		Attributes += food
	}
	food.Attributes = Attributes
	food.DelFlag = 0
	food.Rating = 0 //初始评分为零
	isSuccess := cc.Service.SaveFood(food)
	if !isSuccess {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_ERROR_FOODADD,
				"message": utils.Recode2Text(utils.RESPMSG_ERROR_FOODADD),
			},
		}
	}
	return mvc.Response{
		Object: map[string]interface{}{
			"status":  utils.RECODE_OK,
			"success": utils.Recode2Text(utils.RESPMSG_SUCCESS_FOODADD),
		},
	}

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
