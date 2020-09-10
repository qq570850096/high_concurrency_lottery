package tickets

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"math/rand"
	"time"
)

type lotterController struct {
	Ctx iris.Context
}

func NewApp() *iris.Application {
	app := iris.New()
	mvc.New(app.Party("/")).Handle(&lotterController{})
	return app
}

func main() {
	app := NewApp()

	app.Run(iris.Addr(":8080"))
}

// 即开即得型地址
func (c *lotterController) Get() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Int31n(10)
	prize := ""
	switch {
	case code == 1 :
		prize = "一等奖"
	case code >= 2 && code <= 3:
		prize = "二等奖"
	case code >= 4 && code <= 6:
		prize = "三等奖"
	default:
		return fmt.Sprintf("尾号为1获得一等奖<br/>"+
			"尾号为2或3获得二等奖<br/>"+
			"尾号为4/5/6获得三等奖<br/>"+
			"code=%d<br/>"+
			"很遗憾没有获奖",code)
	}

	return fmt.Sprintf("尾号为1获得一等奖<br/>"+
		"尾号为2或3获得二等奖<br/>"+
		"尾号为4/5/6获得三等奖<br/>"+
		"code=%d<br/>"+
		"恭喜你获得:%s",code,prize)
}

// 双色球自选型
func (c *lotterController) GetPrize() string {
	rand.Seed(time.Now().UnixNano())
	// 先选6个红球， 1-33
	var prize []int
	prize = append(prize,Knuth(33,6)...)
	// 最后一位蓝球，区间是1-16
	prize = append(prize,Knuth(16,1)...)
	return fmt.Sprintf("今日开奖号码是： %v",prize)
}
// 根据n个数字生成m个中奖序列，要求其等概率
func Knuth(n,m int) []int {
	ret := make([]int,m)
	from := make([]int,n)
	for i := range from{
		from[i] = i+1
	}

	for i := 0; i < m; i++{
		// 先生成一个在[i..n)的随机数，然后与i位数字交换
		x := rand.Int()%(n-i)+i
		from[x],from[i] = from[i],from[x]
		ret[i] = from[i]
	}
	return ret
}