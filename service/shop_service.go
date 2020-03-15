package service

import (
	"CmsProject/model"
	"github.com/go-xorm/xorm"
)

/**
 * 商店Shop的服务
 */
type ShopService interface {
	//查询商店总数，并返回
	GetShopCount() (int64, error)
	GetShopList(offset, limit int) []model.Shop
}

type shopService struct {
	engine *xorm.Engine
}

/**
 * 新实例化一个商店模块服务对象结构体
 */
func NewShopService(db *xorm.Engine) ShopService {
	return &shopService{engine: db}
}

/**
 * 查询商店的总数然后返回
 */
func (ss *shopService) GetShopCount() (int64, error) {
	result, err := ss.engine.Where("dele=?", 0).Count(new(model.Shop))
	if err != nil {
		panic(err.Error())
		return 0, err
	}
	return result, nil
}

/**
 * 获取到商铺列表信息
 */
func (ss *shopService) GetShopList(offset, limit int) []model.Shop {
	var shopList []model.Shop
	ss.engine.Where("dele=?", 0).Limit(limit, offset).Find(&shopList)
	return shopList
}
