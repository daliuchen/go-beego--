package controllers

import (
	"beegoTest/models"
	"beegoTest/util"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type ShopCOntroller struct {
	//method继承，匿名字段，
	// mainController 拥有的了 beego.controller里面的所有方法和属性。
	beego.Controller
}

func (this *ShopCOntroller) Get() {
	o := orm.NewOrm()
	var shops []models.ShopDo
	o.Raw("select id ,name,price,count from t_shop").QueryRows(&shops)
	this.TplName = "shop/index.html"
	this.Data["shopList"] = shops
}
func (this *ShopCOntroller) Buy() {
	id := this.Ctx.Input.Param(":id")
	shop := new(models.ShopDo)
	o := orm.NewOrm()
	o.Raw("select id ,name,price,count from t_shop where id = ?", id).QueryRow(&shop)
	shop.Count = 1
	logger.Println("购买的商品信息为:{}", shop)
	util.AdjustCar(shop)
	//key 必须是json
	//但是有个问题，在json导出的时候 标签不起作用
	result := models.ApiResult{200, "", nil}
	this.Data["json"] = result
	//将data里面的数据 json返回
	this.ServeJSON()
}
func (this *ShopCOntroller) End() {
	//从redis中获取列表展示给前端就好了.
	list := util.GetShopCarList()
	var sum float32 = 0
	for _, value := range *list {
		count := float32(value.Count)
		price := value.Price
		sum += count * price
	}
	this.TplName = "shop/end.html"
	this.Data["shopList"] = list
	this.Data["sum"] = sum
}
func (this *ShopCOntroller) GiveMoney() {
	//从redis中获取列表展示给前端就好了.
	list := util.GetShopCarList()
	fmt.Println(list)
	this.TplName = "shop/giveMoney.html"
	this.Data["orderNo"] = util.SnowFlakeUtil.NextId()
}
