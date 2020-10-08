package routers

import (
	"beegoTest/controllers"
	"github.com/astaxie/beego"
)

func init() {
	//采用beego默认的restful来处理，get请求对应Get方法
	beego.Router("/", &controllers.LoginController{})
	//自定义方法来处理 login请求
	beego.Router("/login", &controllers.LoginController{}, "post:Login")
	beego.Router("/logOut", &controllers.LoginController{}, "get:LogOut")
	//自定义方法来处理 login请求 对应get方法
	beego.Router("/list", &controllers.ShopCOntroller{})
	beego.Router("/buy/:id", &controllers.ShopCOntroller{}, "get:Buy")
	beego.Router("/end", &controllers.ShopCOntroller{}, "get:End")
	beego.Router("/giveMoney", &controllers.ShopCOntroller{}, "get:GiveMoney")
}
