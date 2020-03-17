package service

import (
	"CmsProject/model"
	"github.com/go-xorm/xorm"
)

type FoodService interface {
	GetFoodCount() (int64, error)
	GetFoodList(offset, limit int) ([]model.Food, error)
}

type foodService struct {
	engine *xorm.Engine
}

func NewFoodService(db *xorm.Engine) FoodService {
	return &foodService{engine: db}
}
func (fs *foodService) GetFoodCount() (int64, error) {
	count, err := fs.engine.Count(new(model.Food))
	return count, err
}

func (fs *foodService) GetFoodList(offset, limit int) ([]model.Food, error) {
	var foodList []model.Food
	err := fs.engine.Where("del_flag=?", 0).Limit(limit, offset).Find(&foodList)
	return foodList, err
}
