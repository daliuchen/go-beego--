package models

/**
用户实体类
*/
type UserDo struct {
	Id       int    `id"`
	UserName string `form:"userName"`
	Password string `form:"password"`
	Role     string
}

func (this UserDo) String() string {
	return this.UserName + "\t" + this.Password
}
func (u *UserDo) TableName() string {
	return "t_customer"
}

/**
商品实体类
*/
type ShopDo struct {
	Id    int     `id`
	Name  string  `name`
	Price float32 `price`
	Count int     `count`
}

//公共返回实体类
type ApiResult struct {
	Code int         `code`
	Msg  string      `msg`
	Data interface{} `data`
}
