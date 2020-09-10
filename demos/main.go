package main

import (
	"log"
	"math/rand"
	"os"
	"sort"
	"sync"
	"time"
)

var logger *log.Logger

// 记录中奖序列全局变量
var prize []int

// 随机种子
var randseed *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func main() {
	InitLog()
	var once sync.Once
	for i := 0; i < 10; i++ {
		once.Do(LuckyNum)
	}
	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		wg.Add(100)
		go BuyIt100Times(&wg)
	}
	wg.Wait()
}

func BuyIt100Times(wg *sync.WaitGroup) {
	for i := 0; i < 100; i++ {
		Compare(Choose(), wg)
	}
}

func InitLog() {
	f, _ := os.Create("./lottery_demo.log")
	logger = log.New(f, "", log.Ldate|log.Lmicroseconds)
}

func LuckyNum() {
	// 先选6个红球， 1-33
	prize = append(prize, KnuthNum(33, 6)...)
	// 最后一位蓝球，区间是1-16
	prize = append(prize, KnuthNum(16, 1)...)
	// 使用等于号给切片赋值，是引用传递，也就是说我们修改 sort2 同样也会修改 prize 的前六位元素。
	sort2 := prize[:6]
	sort.Ints(sort2)
	logger.Printf("今日开奖号码是： %v\n", prize)
}

// 根据n个数字生成m个中奖序列，要求其等概率
func KnuthNum(n, m int) []int {
	ret := make([]int, m)
	from := make([]int, n)
	for i := range from {
		from[i] = i + 1
	}

	for i := 0; i < m; i++ {
		// 先生成一个在[i..n)的随机数，然后与i位数字交换
		x := randseed.Int()%(n-i) + i
		from[x], from[i] = from[i], from[x]
		ret[i] = from[i]
	}
	return ret
}

// 即开即得型地址
func Choose() []int {
	choose := make([]int, 0)
	// 先选6个红球， 1-33
	choose = append(choose, KnuthNum(33, 6)...)
	// 最后一位蓝球，区间是1-16
	choose = append(choose, KnuthNum(16, 1)...)
	sort2 := choose[:6]
	sort.Ints(sort2)
	return choose
}

func Compare(choose []int, wg *sync.WaitGroup) bool {
	defer wg.Done()
	for i, v := range choose {
		if v != prize[i] {
			logger.Println("没有中奖。")
			return false
		}
	}
	logger.Println("中大奖了！")
	return true
}
