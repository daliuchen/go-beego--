package controllers

import (
	"beegoTest/models"
	"beegoTest/util"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"log"
	"time"
)

/*
定义一个全局的日志
*/
var logger = logs.GetLogger()

type LoginController struct {
	//method继承，匿名字段，
	// mainController 拥有的了 beego.controller里面的所有方法和属性。
	beego.Controller
}

//采用默认的restful方法，get请求对应get方法
func (this *LoginController) Get() {
	/**
	- 使用模版
		如果用户不设置该参数，那么默认会去到模板目录的 Controller/<方法名>.tpl 默认的模版目录可以改变，在配置文件里面配置
	- 不是用模版
		直接用 this.Ctx.WriteString 输出字符串，
	*/
	time := time.Now()
	//指定模版路径和文件名
	this.TplName = "login.html"
	//设置数据
	this.Data["currentTime"] = time.Format("2006-01-02 15:04:05")
	this.Data["message"] = ""
}

//自定义登录处理方法 参数得是 context.Context
func (this *LoginController) Login() {
	/*
		第一种，这中是很简单的，和Spring中的param一样。一个一个获取。
	*/
	//userName := c.GetString("userName")
	//password := c.GetString("password")
	//logger.Println(userName,password)

	u := models.UserDo{}
	if err := this.ParseForm(&u); err != nil {
		logs.Error("解析用户登录参数错误", err)
		this.TplName = "exception/500.html"
		return
	}
	//验证
	valid := validation.Validation{}
	valid.Required(u.UserName, "UserName")
	valid.MaxSize(u.UserName, 6, "UserName")
	if valid.HasErrors() {
		// 如果有错误信息，证明验证没通过
		// 打印错误信息
		for _, err := range valid.Errors {
			log.Println(err.Key, err.Message)
		}
		this.TplName = "login.html"
		this.Data["message"] = "用户名或密码错误"
		return
	}

	if this.validataUser(u) == false {
		logger.Println("用户名或密码错误")
		this.TplName = "login.html"
		this.Data["message"] = "用户名或密码错误"
		return
	}
	if u.UserName == "admin" {
		this.Data["isAdmin"] = true
	}
	logger.Println("登录成功", u)
	logger.Println("创建购物车")
	this.Redirect("/list", 302)
}
func (this *LoginController) validataUser(user models.UserDo) bool {
	o := orm.NewOrm()
	var count int
	err := o.Raw("select count(1) from t_customer where name = ? and password = ?", user.UserName, user.Password).QueryRow(&count)
	if err != nil {
		logger.Println(err)
	}
	if count == 1 {
		return true
	} else {
		return false
	}
}
func (this *LoginController) LogOut() {
	util.ReleaseRedisCar()
	this.CruSession = nil
	this.TplName = "login.html"
}
