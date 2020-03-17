package main

import (
	"CmsProject/config"
	"CmsProject/controller"
	"CmsProject/datasource"
	"CmsProject/model"
	"CmsProject/service"
	"CmsProject/utils"
	"encoding/json"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	app := newApp()

	//应用App设置
	ConfigAction(app)

	//路由设置
	mvcHandle(app)

	config := config.InitConfig()
	addr := ":" + config.Port
	app.Run(
		iris.Addr(addr), //在端口8080进行监听
		iris.WithoutServerError(iris.ErrServerClosed), //无服务错误提示
		iris.WithOptimizations,                        //对json数据序列化更快的配置
	)
}

//构建App
func newApp() *iris.Application {
	app := iris.New()

	//设置日志级别  开发阶段为debug
	app.Logger().SetLevel("debug")

	//注册静态资源
	app.HandleDir("/static", "./static")
	app.HandleDir("/manage/static", "./static")
	app.HandleDir("/img", "./uploads")

	//注册视图文件
	app.RegisterView(iris.HTML("./static", ".html"))
	app.Get("/", func(context context.Context) {
		context.View("index.html")
	})

	return app
}

/**
 * MVC 架构模式处理
 */
func mvcHandle(app *iris.Application) {

	//启用session
	sessManager := sessions.New(sessions.Config{
		Cookie:  "sessioncookie",
		Expires: 24 * time.Hour,
	})

	//获取redis实例
	redis := datasource.NewRedis()
	//设置session的同步位置为redis
	sessManager.UseDatabase(redis)

	//实例化mysql数据库引擎
	engine := datasource.NewMysqlEngine()

	//管理员模块功能
	adminService := service.NewAdminService(engine)
	admin := mvc.New(app.Party("/admin"))
	admin.Register(adminService, sessManager.Start)
	admin.Handle(new(controller.AdminController))

	//统计功能模块
	statisService := service.NewStatisService(engine)
	statis := mvc.New(app.Party("/statis/{model}/{date}/"))
	statis.Register(statisService, sessManager.Start)
	statis.Handle(new(controller.StatisController))

	//用户功能模块
	userService := service.NewUserService(engine)
	user := mvc.New(app.Party("/v1/users/"))
	user.Register(userService, sessManager.Start)
	user.Handle(new(controller.UserController))

	//订单模块
	orderService := service.NewOrderService(engine)
	order := mvc.New(app.Party("/bos/orders/"))
	order.Register(orderService, sessManager.Start)
	order.Handle(new(controller.OrderController))

	//商铺模块
	shopService := service.NewShopService(engine)
	shop := mvc.New(app.Party("/shopping/restaurants/"))
	shop.Register(shopService, sessManager.Start)
	shop.Handle(new(controller.ShopController))

	//添加食品类别
	categoryService := service.NewCategoryService(engine)
	category := mvc.New(app.Party("/shopping/"))
	category.Register(categoryService, sessManager.Start)
	category.Handle(new(controller.CategoryController))

	foodService := service.NewFoodService(engine)
	food := mvc.New(app.Party("/shopping/v2/foods/"))
	food.Register(foodService)
	food.Handle(new(controller.FoodController))

	//获取用户详细信息
	app.Get("/v1/user/{user_name}", func(context context.Context) {
		userName := context.Params().Get("user_name")
		var user model.User
		_, err := engine.Where("user_name=?", userName).Get(&user)
		if err != nil {
			context.JSON(iris.Map{
				"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_ERROR_USERINFO,
				"message": utils.Recode2Text(utils.RESPMSG_ERROR_USERINFO),
			})
		} else {
			context.JSON(user)
		}
	})

	//获取地址信息
	app.Get("/v1/addresse/{address_id}", func(context context.Context) {
		addressId := context.Params().Get("address_id")
		AddressId, err := strconv.Atoi(addressId)
		if err != nil {
			context.JSON(iris.Map{
				"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_ERROR_ORDERINFO,
				"message": utils.Recode2Text(utils.RESPMSG_ERROR_ORDERINFO),
			})
		}

		var address model.Address
		_, err = engine.Where("address_id=?", AddressId).Get(&address)
		if err != nil {
			context.JSON(iris.Map{
				"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_ERROR_ORDERINFO,
				"message": utils.Recode2Text(utils.RESPMSG_ERROR_ORDERINFO),
			})
		} else {

			//查询数据成功
			context.JSON(address)
		}
	})

	//文件上传
	app.Post("/admin/update/avatar/{adminId}", func(context context.Context) {
		adminId := context.Params().Get("adminId")
		adminID, _ := strconv.Atoi(adminId)
		file, info, err := context.FormFile("file")
		if err != nil {
			context.JSON(iris.Map{
				"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_ERROR_PICTUREADD,
				"failure": utils.Recode2Text(utils.RESPMSG_ERROR_PICTUREADD),
			})
			return
		}
		fname := info.Filename
		out, err := os.OpenFile("./uploads/"+fname, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			context.JSON(iris.Map{
				"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_ERROR_PICTUREADD,
				"failure": utils.Recode2Text(utils.RESPMSG_ERROR_PICTUREADD),
			})
			return
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			context.JSON(iris.Map{
				"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_ERROR_PICTUREADD,
				"failure": utils.Recode2Text(utils.RESPMSG_ERROR_PICTUREADD),
			})
			return
		}

		adminService.SaveAvatarImg(int64(adminID), fname)

		context.JSON(iris.Map{
			"status":     utils.RECODE_OK,
			"image_path": fname,
		})
	})
	app.Get("/shopping/restaurant/{shopping_id}", func(context context.Context) {
		shopping_id := context.Params().Get("shopping_id")
		ShoppingId, _ := strconv.Atoi(shopping_id)

		var shop model.Shop
		_, err := engine.Where("shop_id=?", ShoppingId).Get(&shop)
		if err != nil {
			panic(err.Error())
			return
		}
		context.JSON(shop)
	})

	app.Post("/v1/addimg/{model}}", func(context context.Context) {
		model := context.Params().Get("model")
		file, info, err := context.FormFile("file")
		if err != nil {
			context.JSON(iris.Map{
				"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_ERROR_PICTUREADD,
				"failure": utils.Recode2Text(utils.RESPMSG_ERROR_PICTUREADD),
			})
			return
		}
		defer file.Close()
		fname := info.Filename
		isExist, err := utils.PathExists("./uploads/" + model)
		if !isExist {
			err := os.Mkdir("./uploads/"+model, 0777)
			if err != nil {
				context.JSON(iris.Map{
					"status":  utils.RECODE_FAIL,
					"type":    utils.RESPMSG_ERROR_PICTUREADD,
					"failure": utils.Recode2Text(utils.RESPMSG_ERROR_PICTUREADD),
				})
				return
			}
		}
		out, err := os.OpenFile("./uploads/"+model+"/"+fname, os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			context.JSON(iris.Map{
				"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_ERROR_PICTUREADD,
				"failure": utils.Recode2Text(utils.RESPMSG_ERROR_PICTUREADD),
			})
			return
		}
		_, err = io.Copy(out, file)
		if err != nil {
			context.JSON(iris.Map{
				"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_ERROR_PICTUREADD,
				"failure": utils.Recode2Text(utils.RESPMSG_ERROR_PICTUREADD),
			})
			return
		}
		context.JSON(iris.Map{
			"status":     utils.RECODE_OK,
			"image_path": fname,
		})
	})

	//地址Poi检索
	app.Get("/v1/pois/", func(context context.Context) {
		path := context.Request().URL.String()
		rs, err := http.Get("https://elm.cangdu.org" + path)
		if err != nil {
			context.JSON(iris.Map{
				"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_ERROR_SEARCHADDRESS,
				"message": utils.Recode2Text(utils.RESPMSG_ERROR_SEARCHADDRESS),
			})
			return
		}
		//请求成功
		body, err := ioutil.ReadAll(rs.Body)
		if err != nil {
			panic(err.Error())
			return
		}
		var searchList []*model.PoiSearch
		json.Unmarshal(body, &searchList)
		context.JSON(&searchList)
	})

}

/**
 * 项目设置
 */
func ConfigAction(app *iris.Application) {

	//配置 字符编码
	app.Configure(iris.WithConfiguration(iris.Configuration{
		Charset: "UTF-8",
	}))

	//错误配置
	//未发现错误
	app.OnErrorCode(iris.StatusNotFound, func(context context.Context) {
		context.JSON(map[string]interface{}{
			"errmsg": iris.StatusNotFound,
			"msg":    " not found ",
			"data":   iris.Map{},
		})
	})

	//错误配置
	//服务器内部错误
	app.OnErrorCode(iris.StatusInternalServerError, func(context context.Context) {
		context.JSON(map[string]interface{}{
			"errmsg": iris.StatusInternalServerError,
			"msg":    " interal error  ",
			"data":   iris.Map{},
		})
	})
}
