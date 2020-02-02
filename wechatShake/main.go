package wechatShake

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)
var mu sync.Mutex
// 微信摇一摇基础功能/lucky 一个抽奖的接口
const (
	//奖品类型，枚举值iota会从0开始
	GiftTypeCoin      = iota // 虚拟币
	GiftTypeCoupon           // 不同的券
	GiftTypeCouponFix        // 相同的券
	GiftTypeRealSmall        // 小奖
	GiftTypeRealLarge        // 大奖
)

type gift struct {
	id       int
	name     string
	pic      string // 奖品图片
	link     string // 奖品链接
	gtype    int
	desc     string   // 奖品描述
	dataList []string // 奖品数据集合（比如不同优惠券的编码）
	total    int
	left     int
	inuse    bool
	rate     int // 中奖概率，万分之n
	rateMin  int // 中奖编码最小值
	rateMax  int // 中奖编码最大值
}

// 最大中奖号码
const rateMax = 10000

var logger *log.Logger

// 奖品列表
var giftlist []*gift

type lotteryController struct {
	Ctx iris.Context
}

func InitLog() {
	f, _ := os.Create("./lottery_demo.log")
	logger = log.New(f, "", log.Ldate|log.Lmicroseconds)
}

func InitGift() {
	giftlist = make([]*gift, 5)
	g1 := &gift{
		id:       0,
		name:     "Mate30Pro 5G",
		pic:      "",
		link:     "",
		gtype:    GiftTypeRealLarge,
		desc:     "",
		dataList: nil,
		total:    20000,
		left:     20000,
		inuse:    true,
		rate:     10000,
		rateMin:  0,
		rateMax:  0,
	}
	giftlist[0] = g1
	g2 := &gift{
		id:       0,
		name:     "充电宝",
		pic:      "",
		link:     "",
		gtype:    GiftTypeRealSmall,
		desc:     "",
		dataList: nil,
		total:    5,
		left:     5,
		inuse:    false,
		rate:     10,
		rateMin:  0,
		rateMax:  0,
	}
	giftlist[1] = g2
	g3 := &gift{
		id:       0,
		name:     "优惠券200-50",
		pic:      "",
		link:     "",
		gtype:    GiftTypeCouponFix,
		desc:     "coupon-2020",
		dataList: nil,
		total:    50,
		left:     50,
		inuse:    false,
		rate:     500,
		rateMin:  0,
		rateMax:  0,
	}
	giftlist[2] = g3
	g4 := &gift{
		id:       0,
		name:     "折价50元优惠券",
		pic:      "",
		link:     "",
		gtype:    GiftTypeCoupon,
		desc:     "",
		dataList: []string{"c01", "c02", "c03", "c04", "c05"},
		total:    5,
		left:     5,
		inuse:    false,
		rate:     500,
		rateMin:  0,
		rateMax:  0,
	}
	giftlist[3] = g4
	g5 := &gift{
		id:       0,
		name:     "金币",
		pic:      "",
		link:     "",
		gtype:    GiftTypeCoin,
		desc:     "10金币",
		dataList: nil,
		total:    50,
		left:     50,
		inuse:    false,
		rate:     5000,
		rateMin:  0,
		rateMax:  0,
	}
	giftlist[4] = g5
	// 数据整理， 中奖区间数据
	rateStart := 0
	for _, data := range giftlist {
		if !data.inuse {
			continue
		}
		data.rateMin = rateStart
		data.rateMax = rateStart + data.rate

		if data.rateMax >= rateMax {
			data.rateMax = rateMax
			rateStart = 0
		} else {
			rateStart += data.rate
		}
	}
}

func NewApp() *iris.Application {
	app := iris.New()

	mvc.New(app.Party("/")).Handle(&lotteryController{})
	InitGift()
	InitLog()
	return app
}

func main() {
	app := NewApp()

	app.Run(iris.Addr(":8080"))
}

// 奖品数量信息
func (c *lotteryController) Get() string {
	count := 0
	total := 0
	for _, v := range giftlist {
		if v.inuse && (v.total == 0 ||
			(v.total > 0 && v.left > 0)) {
			count++
			total += v.left
		}
	}
	return fmt.Sprintf("当前有效奖品种类数量:%d,限量奖品总量：%d\n", count, total)
}

// GET http://localhost:8080/lucky
func (c *lotteryController) GetLucky() map[string]interface{} {
	mu.Lock()
	defer mu.Unlock()
	code := luckyCode()
	ok := false
	result := make(map[string]interface{})
	result["success"] = ok
	for _, data := range giftlist {
		if !data.inuse || (data.total > 0 && data.left <= 0) {
			continue
		}
		if data.rateMin <= int(code) && data.rateMax > int(code) {
			// 中奖了，抽奖编码在奖品中奖编码范围内
			sendData := ""
			switch data.gtype {
			case GiftTypeCoin:
				ok, sendData = sendCoin(data)
			case GiftTypeCoupon:
				ok, sendData = sendCoupon(data)
			case GiftTypeCouponFix:
				ok, sendData = sendCouponFix(data)
			case GiftTypeRealSmall:
				ok, sendData = sendRealSmall(data)
			case GiftTypeRealLarge:
				ok, sendData = sendRealLarge(data)
			}
			if ok {
				// 中奖后，成功得到奖品（发奖成功）
				// 生成中奖纪录
				saveLuckyData(code, data.id, data.name, data.link, sendData, data.left)
				result["success"] = ok
				result["id"] = data.id
				result["name"] = data.name
				result["link"] = data.link
				result["data"] = sendData
				break
			}
		}
	}

	return result
}

// 抽奖编码
func luckyCode() int32 {
	seed := time.Now().UnixNano()                                 // rand内部运算的随机数
	code := rand.New(rand.NewSource(seed)).Int31n(int32(rateMax)) // rand计算得到的随机数
	return code
}

// 发奖，虚拟币
func sendCoin(data *gift) (bool, string) {
	if data.total == 0 {
		// 数量无限
		return true, data.desc
	} else if data.left > 0 {
		// 还有剩余
		data.left = data.left - 1
		return true, data.desc
	} else {
		return false, "奖品已发完"
	}
}

// 发奖，优惠券（不同值）
func sendCoupon(data *gift) (bool, string) {
	if data.left > 0 {
		// 还有剩余的奖品
		left := data.left - 1
		data.left = left
		return true, data.dataList[left]
	} else {
		return false, "奖品已发完"
	}
}

// 发奖，优惠券（固定值）
func sendCouponFix(data *gift) (bool, string) {
	if data.total == 0 {
		// 数量无限
		return true, data.desc
	} else if data.left > 0 {
		data.left = data.left - 1
		return true, data.desc
	} else {
		return false, "奖品已发完"
	}
}

// 发奖，实物小
func sendRealSmall(data *gift) (bool, string) {
	if data.total == 0 {
		// 数量无限
		return true, data.desc
	} else if data.left > 0 {
		data.left = data.left - 1
		return true, data.desc
	} else {
		return false, "奖品已发完"
	}
}

// 发奖，实物大
func sendRealLarge(data *gift) (bool, string) {
	if data.total == 0 {
		// 数量无限
		return true, data.desc
	} else if data.left > 0 {
		data.left--
		return true, data.desc
	} else {
		return false, "奖品已发完"
	}
}

// 记录用户的获奖记录
func saveLuckyData(code int32, id int, name, link, sendData string, left int) {
	logger.Printf("lucky, code=%d, gift=%d, name=%s, link=%s, data=%s, left=%d ", code, id, name, link, sendData, left)
}
