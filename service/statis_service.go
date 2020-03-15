package service

import (
	"CmsProject/model"
	"fmt"
	"github.com/go-xorm/xorm"
	"math/rand"
	"time"
)

/**
 * 统计功能模块接口标准
 */
type StatisService interface {
	//查询某一天的用户的增长数量
	GetUserDailyCount(date string) int64

	//查询某一天的订单的增长数量
	GetOrderDailyCount(date string) int64

	//查询某一天的管理员的增长数量
	GetAdminDailyCount(date string) int64
}

/**
 * 统计功能服务实现结构体
 */
type statisService struct {
	engine *xorm.Engine
}

/**
 * 新建统计模块功能服务对象
 */
func NewStatisService(db *xorm.Engine) StatisService {
	return &statisService{engine: db}
}

/**
 * 查询某一日管理员的增长数量
 */
func (ss *statisService) GetAdminDailyCount(date string) int64 {
	//查询日期date格式解析
	startDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		panic(err.Error())
		return 0
	}
	endDate := startDate.AddDate(0, 0, 1)
	result, err := ss.engine.Where("create_time between ? and ? and status=0", startDate.Format("2006-01-02 15:04:05"),
		endDate.Format("2006-01-02 15:04:05")).Count(model.Admin{})
	if err != nil {
		panic(err.Error())
		return 0
	}
	fmt.Println(date, "\t管理员总数:", result)
	return int64(rand.Intn(100))
}

/**
 * 查询某一日订单的单日增长数量
 */
func (ss *statisService) GetOrderDailyCount(date string) int64 {
	startDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		panic(err.Error())
		return 0
	}
	endDate := startDate.AddDate(0, 0, 1)
	result, err := ss.engine.Where("time between ? and ? and del_flag=0", startDate.Format("2006-01-02 15:04:05"),
		endDate.Format("2006-01-02 15:04:05")).Count(model.UserOrder{})
	if err != nil {
		panic(err.Error())
		return 0
	}
	fmt.Println(date, "\t订单总数:", result)
	return int64(rand.Intn(100))
}

/**
 * 查询某一日用户的单日增长数量
 */
func (ss *statisService) GetUserDailyCount(date string) int64 {
	startDate, err := time.Parse("2006-01-02", date)

	if err != nil {
		panic(err.Error())
		return 0
	}
	endDate := startDate.AddDate(0, 0, 1)
	result, err := ss.engine.Where("register_time between ? and ? and del_flag=0", startDate.Format("2006-01-02 15:04:05"),
		endDate.Format("2006-01-02 15:04:05")).Count(model.User{})
	if err != nil {
		panic(err.Error())
		return 0
	}
	fmt.Println(date, "\t用户总数:", result)
	return int64(rand.Intn(100))
}
