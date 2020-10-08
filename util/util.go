package util

import (
	"beegoTest/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/bitly/go-simplejson"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"sync"
	"time"
)

var RedisConnection *redis.Conn
var SnowFlakeUtil *SnowFlake

func getConfigValueByName(key string) string {
	return beego.AppConfig.String("key")
}
func init() {
	c, err := redis.Dial("tcp", "127.0.0.1:6378")
	RedisConnection = &c
	if err != nil {
		fmt.Println("connect redis error :", err)
		return
	}

	//初始化SnowFlake
	SnowFlakeUtil = new(SnowFlake)
	SnowFlakeUtil.Init(1, 1, 1)
}
func AdjustCar(shop *models.ShopDo) {
	carStr, _ := (*RedisConnection).Do("get", "car")
	fmt.Println(carStr)
	if carStr == nil {
		var cars []models.ShopDo
		cars = append(cars, *shop)
		b, _ := json.Marshal(cars)
		(*RedisConnection).Do("SET", "car", b)
	} else {
		isRepeat := false
		var cars1 []models.ShopDo
		js, _ := simplejson.NewJson(carStr.([]uint8))
		array, _ := js.Array()
		for _, i2 := range array {
			aa := i2.(map[string]interface{})
			b, _ := strconv.Atoi(aa["Id"].(json.Number).String())
			i3 := aa["Name"]
			float1212, _ := strconv.ParseFloat(aa["Price"].(json.Number).String(), 32)
			i5 := aa["Count"].(json.Number).String()
			atoi, _ := strconv.Atoi(i5)
			if b == shop.Id {
				atoi += 1
				isRepeat = true
			}
			do := models.ShopDo{Id: b, Name: i3.(string), Price: float32(float1212), Count: atoi}
			cars1 = append(cars1, do)
		}
		if isRepeat == false {
			cars1 = append(cars1, *shop)
		}
		b, _ := json.Marshal(cars1)
		fmt.Println("购物车所有商品为", cars1)
		(*RedisConnection).Do("SET", "car", b)
	}
}
func ReleaseRedisCar() {
	fmt.Println("销毁购物车")
	(*RedisConnection).Do("DEL", "car")
}
func GetShopCarList() *[]models.ShopDo {
	carStr, _ := (*RedisConnection).Do("get", "car")
	var cars1 []models.ShopDo
	js, _ := simplejson.NewJson(carStr.([]uint8))
	array, _ := js.Array()
	for _, i2 := range array {
		aa := i2.(map[string]interface{})
		b, _ := strconv.Atoi(aa["Id"].(json.Number).String())
		i3 := aa["Name"]
		float1212, _ := strconv.ParseFloat(aa["Price"].(json.Number).String(), 32)
		i5 := aa["Count"].(json.Number).String()
		atoi, _ := strconv.Atoi(i5)
		do := models.ShopDo{Id: b, Name: i3.(string), Price: float32(float1212), Count: atoi}
		cars1 = append(cars1, do)
	}
	return &cars1
}

type SnowFlake struct {
	//因为二进制里第一个 bit 为如果是 1，那么都是负数，但是我们生成的 id 都是正数，所以第一个 bit 统一都是 0。

	//机器ID  2进制5位  32位减掉1位 31个
	workerId uint64
	//机房ID 2进制5位  32位减掉1位 31个
	datacenterId uint64
	//代表一毫秒内生成的多个id的最新序号  12位 4096 -1 = 4095 个
	sequence uint64
	//设置一个时间初始值    2^41 - 1   差不多可以用69年
	twepoch uint64
	//5位的机器id
	workerIdBits uint64
	//5位的机房id
	datacenterIdBits uint64
	//每毫秒内产生的id数 2 的 12次方
	sequenceBits uint64
	// 这个是二进制运算，就是5 bit最多只能有31个数字，也就是说机器id最多只能是32以内
	maxWorkerId uint64
	// 这个是一个意思，就是5 bit最多只能有31个数字，机房id最多只能是32以内
	maxDatacenterId uint64

	workerIdShift      uint64
	datacenterIdShift  uint64
	timestampLeftShift uint64
	sequenceMask       uint64
	//记录产生时间毫秒数，判断是否是同1毫秒
	lastTimestamp uint64
	//互斥锁
	lock sync.Mutex
}

func (this *SnowFlake) Init(workerId uint64, datacenterId uint64, sequence uint64) {
	// 检查机房id和机器id是否超过31 不能小于0
	this.twepoch = 1585644268888
	this.workerIdBits = 5
	this.datacenterIdBits = 5
	this.sequenceBits = 12
	this.maxWorkerId = uint64(1 ^ (1 << this.workerIdBits))
	this.maxDatacenterId = uint64(1 ^ (1 << this.datacenterIdBits))
	this.workerIdShift = this.sequenceBits
	this.datacenterIdShift = this.sequenceBits + this.workerIdBits
	this.timestampLeftShift = this.sequenceBits + this.workerIdBits + this.datacenterIdBits
	this.sequenceMask = uint64(1 ^ (1 << this.sequenceBits))
	this.lastTimestamp = uint64(1)
	if workerId > this.maxWorkerId || workerId < 0 {
		panic(errors.New("检查机房id和机器id是否超过31 不能小于0"))
	}
	if datacenterId > this.maxDatacenterId || datacenterId < 0 {
		panic(errors.New("检查机房id和机器id是否超过31 不能小于0"))
	}
	this.workerId = workerId
	this.datacenterId = datacenterId
	this.sequence = sequence
}
func (this *SnowFlake) NextId() uint64 {
	this.lock.Lock()
	defer this.lock.Unlock()
	// 这儿就是获取当前时间戳，单位是毫秒
	timestamp := uint64(time.Now().Unix())
	if timestamp < this.lastTimestamp {

	}

	// 下面是说假设在同一个毫秒内，又发送了一个请求生成一个id
	// 这个时候就得把seqence序号给递增1，最多就是4096
	if this.lastTimestamp == timestamp {

		// 这个意思是说一个毫秒内最多只能有4096个数字，无论你传递多少进来，
		//这个位运算保证始终就是在4096这个范围内，避免你自己传递个sequence超过了4096这个范围
		this.sequence = (this.sequence + 1) & this.sequenceMask
		//当某一毫秒的时间，产生的id数 超过4095，系统会进入等待，直到下一毫秒，系统继续产生ID
		if this.sequence == 0 {
			timestamp = this.tilNextMillis(this.lastTimestamp)
		}

	} else {
		this.sequence = 0
	}
	// 这儿记录一下最近一次生成id的时间戳，单位是毫秒
	this.lastTimestamp = timestamp
	// 这儿就是最核心的二进制位运算操作，生成一个64bit的id
	// 先将当前时间戳左移，放到41 bit那儿；将机房id左移放到5 bit那儿；将机器id左移放到5 bit那儿；将序号放最后12 bit
	// 最后拼接起来成一个64 bit的二进制数字，转换成10进制就是个long型
	return ((timestamp - this.twepoch) << this.timestampLeftShift) |
		(this.datacenterId << this.datacenterIdShift) |
		(this.workerId << this.workerIdShift) | this.sequence
}
func (this *SnowFlake) tilNextMillis(lastTimestamp uint64) uint64 {
	/**
	 * 当某一毫秒的时间，产生的id数 超过4095，系统会进入等待，直到下一毫秒，系统继续产生ID
	 * @param lastTimestamp
	 * @return
	 */

	timestamp := uint64(time.Now().Unix())

	for timestamp <= lastTimestamp {
		timestamp = uint64(time.Now().Unix())
	}
	return timestamp
}
