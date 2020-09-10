package tickets

import (
	"log"
	"math/rand"
	"os"
	"sort"
	"sync"
	"testing"
	"time"
)
// 随机种子
var randseed *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
// 中奖序列
var prize []int
var logger *log.Logger
func TestNewApp(t *testing.T) {
	var once sync.Once
	once.Do(LuckyNum)
	var wg sync.WaitGroup
	for i:=0;i<10;i++ {
		wg.Add(1)
		go Compare(Choose(),&wg)
	}
	wg.Wait()
}

func InitLog(){
	f, _ := os.Create("./lottery_demo.log")
	logger = log.New(f, "", log.Ldate|log.Lmicroseconds)
}

// 即开即得型地址
func Choose() []int {
	choose := make([]int,0)
	// 先选6个红球， 1-33
	choose = append(choose,Knuth(33,6)...)
	// 最后一位蓝球，区间是1-16
	choose = append(choose,Knuth(16,1)...)
	sort2 := choose[:6]
	sort.Ints(sort2)
	return choose
}

func Compare(choose []int,wg *sync.WaitGroup) bool {
	defer wg.Done()
	for i,v := range choose {
		if v != prize[i] {
			logger.Println("没有中奖。")
			return false
		}
	}
	logger.Println("中大奖了！")
	return true
}

// 双色球自选型
func LuckyNum()  {
	// 先选6个红球， 1-33
	prize = append(prize,Knuth(33,6)...)
	// 最后一位蓝球，区间是1-16
	prize = append(prize,Knuth(16,1)...)
	sort2 := prize[:6]
	sort.Ints(sort2)
	logger.Printf("今日开奖号码是： \n%v",prize)
}
// 根据n个数字生成m个中奖序列，要求其等概率
func KnuthNum(n,m int) []int {
	ret := make([]int,m)
	from := make([]int,n)
	for i := range from{
		from[i] = i+1
	}

	for i := 0; i < m; i++{
		// 先生成一个在[i..n)的随机数，然后与i位数字交换
		x := randseed.Int()%(n-i)+i
		from[x],from[i] = from[i],from[x]
		ret[i] = from[i]
	}
	return ret
}