package service

import (
	"CmsProject/model"
	"github.com/go-xorm/xorm"
)

/**
 * 食品种类服务接口
 */
type CategoryService interface {
	AddCategory(category *model.FoodCategory) bool
	GetCategoryByShopId(shopId int64) ([]model.FoodCategory, error)
	SaveFood(food model.Food) bool
	DelFood(foodId int) bool
	SaveShop(shop model.Shop) bool
	DeleteShop(shopId int) bool
}

/**
 * 种类服务实现结构体
 */
type categoryService struct {
	engine *xorm.Engine
}

/**
 * 实例化种类服务:服务器
 */
func NewCategoryService(db *xorm.Engine) CategoryService {
	return &categoryService{engine: db}
}

func (cs *categoryService) DeleteShop(shopId int) bool {
	shop := model.Shop{ShopId: shopId, Dele: 1}
	_, err := cs.engine.Where("shop_id=?", shopId).Update(&shop)
	return err == nil
}

func (cs *categoryService) SaveShop(shop model.Shop) bool {
	_, err := cs.engine.Insert(&shop)
	return err == nil
}

func (cs *categoryService) DelFood(foodId int) bool {
	food := model.Food{Id: foodId, DelFlag: 1}
	_, err := cs.engine.Where("id = ?", foodId).Update(&food)
	return err == nil
}

func (cs *categoryService) SaveFood(food model.Food) bool {
	_, err := cs.engine.Insert(&food)
	return err == nil
}

/**
 * 通过商铺Id获取食品类型
 */
func (cs *categoryService) GetCategoryByShopId(shopId int64) ([]model.FoodCategory, error) {
	var foodCategory []model.FoodCategory
	err := cs.engine.Where("restaurant_id = ?", shopId).Find(&foodCategory)
	if err != nil {
		panic(err.Error())
		return nil, err
	}
	return foodCategory, nil
}

/**
 * 添加食品种类记录
 */
func (cs *categoryService) AddCategory(category *model.FoodCategory) bool {
	_, err := cs.engine.Insert(category)
	if err != nil {
		panic(err.Error())
		return false
	}
	return true
}
