package annualMeeting

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"math/rand"
	"strings"
	"sync"
	"time"
)
var mutex sync.Mutex
var userList []string

type lotteryController struct {
	Ctx iris.Context
}

func newApp() *iris.Application {
	app := iris.New()
	mvc.New(app.Party("/")).Handle(&lotteryController{})
	return app
}

func main() {
	app := newApp()
	userList = make([]string, 0)

	app.Run(iris.Addr(":8080"))
}

func (c *lotteryController)Get() string {
	count := len(userList)
	return fmt.Sprintf("当前总共参与抽奖的用户数为： %d",count)
}
func (c *lotteryController) PostImport() string {
	mutex.Lock()
	defer mutex.Unlock()
	strUsers := c.Ctx.FormValue("users")
	users := strings.Split(strUsers,",")
	count1 := len(userList)
	for _,u := range users {
		u = strings.TrimSpace(u)
		if len(u) > 0 {
			userList = append(userList, u)
		}
	}
	count2 := len(userList)
	return fmt.Sprintf("当前总共参与的人数： %d\n，" +
		"成功导入的用户数： %d\n",count2,count2-count1)
}

func (c *lotteryController) GetLucky() string {
	mutex.Lock()
	defer mutex.Unlock()
	count := len(userList)
	if count > 1 {
		rand.Seed(time.Now().UnixNano())
		index := rand.Int31n(int32(count))
		lucky := userList[index]
		userList = append(userList[0:index], userList[index+1:]...)
		return fmt.Sprintf("当前中奖用户： %s，剩余用户数： %d\n",lucky,count-1)
	} else if count == 1{
		lucky := userList[0]
		userList = append(userList[1:])
		return fmt.Sprintf("当前中奖用户： %s，剩余用户数： %d\n",lucky,count-1)
	} else {
		return fmt.Sprintf("当前没有用户，请先import\n")
	}
}