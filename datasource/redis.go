package datasource

import (
	"CmsProject/config"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions/sessiondb/redis"
)

/**
 * 返回Redis实例
 */
func NewRedis() *redis.Database {
	var database *redis.Database

	cmsConfig := config.InitConfig()
	if cmsConfig != nil {
		rd := cmsConfig.Redis
		database = redis.New(redis.Config{
			Network:   rd.NetWork,
			Addr:      rd.Addr + ":" + rd.Port,
			Password:  rd.Password,
			Database:  "",
			MaxActive: 10,
			Timeout:   redis.DefaultRedisTimeout,
			Prefix:    rd.Prefix,
		})
	} else {
		iris.New().Logger().Info("error")
	}
	return database
}
