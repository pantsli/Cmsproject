package service

import (
	"CmsProject/model"
	"github.com/go-xorm/xorm"
)

/**
 * 订单服务接口
 */
type OrderService interface {
	//获取订单总数量
	GetCount() (int64, error)

	//获取订单列表
	GetOrderList(offset, limit int) []model.OrderDetail
}

/**
 * 订单服务
 */
type orderService struct {
	engine *xorm.Engine
}

/**
 * 实例化OrderService服务对象
 */
func NewOrderService(db *xorm.Engine) OrderService {
	return &orderService{engine: db}
}

/**
 * 获取订单列表
 */
func (os *orderService) GetOrderList(offset, limit int) []model.OrderDetail {
	var orderList []model.OrderDetail

	err := os.engine.Table("user_order").
		Join("INNER", "order_status", " order_status.status_id = user_order.order_status_id   ").
		Join("INNER", "user", " user.id = user_order.user_id").
		Join("INNER", "shop", " shop.shop_id = user_order.shop_id ").
		Join("INNER", "address", " address.address_id = user_order.address_id ").
		Find(&orderList)

	if err != nil {
		panic(err.Error())
		return nil
	}

	return orderList
}

/**
 * 获取订单总数量
 */
func (os *orderService) GetCount() (int64, error) {
	count, err := os.engine.Where(" del_flag = ?", 0).Count(new(model.UserOrder))
	if err != nil {
		return 0, err
	}
	return count, nil
}
