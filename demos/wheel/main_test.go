package wheel

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"
)
var count int = 0

func TestLotteryController_Get(t *testing.T) {
	main()
}

func BenchmarkNewApp(b *testing.B) {
	gof := func(wg *sync.WaitGroup){
		defer wg.Done()
		var resp *http.Response
		var err error
		for i:=0; i<1005; i++{
			//模拟一个get提交请求
			resp, err = http.Get("http://localhost:8080/prize")
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close() //关闭连接
		}
	}
	var wg sync.WaitGroup

	for i:=0 ;i < 10 ; i++ {
		wg.Add(1)
		// 5个协程，应该是10000个数据
		go gof(&wg)
	}
	wg.Wait()
	// 第一次测试后时隔协程并发执行后应该发了一万张票然后50个谢谢参与
	// 实测多发八张票，这个就很tm恐怖，超发了，你自裁罢。
	// 使用mutex锁时，10个协程共用时2856364800 ns
	// 使用原子操作时，十个协程共用时2751644700 ns
}



func TestInitLog(t *testing.T) {
	product := make(chan int,10)
	ctx,cancel := context.WithCancel(context.TODO())
	go Consumer(product, ctx,cancel)
	for i:=0; i < 10 ; i++ {
		go Producter(product,ctx)
	}
	<-ctx.Done()
}

func Producter(ch chan<- int,ctx context.Context) {
	for i := 0; i < 10 ; i++ {
		fmt.Println("生产了产品",count)
		ch <- count
		count++
	}
	// 生产结束了
	if count == 100 {
		close(ch)
	}
	return
}

func Consumer(ch <-chan int,ctx context.Context,cancelFunc context.CancelFunc) {
	for {
		select {
		case v,ok := <-ch:
			// ch还没关闭，也就是生产活动还没结束
			if ok {
				fmt.Printf("消费者消费了产品%d\n ", v)
			} else { // ch 已经关闭了，以后都不会再生产了
				fmt.Println("所有生产的商品都已经消费完了")
				cancelFunc()
				// 想从 for select 中跳出，break是不行的，只能用goto跳出
				goto END
			}
		default:
			fmt.Printf("目前没有产品了，正在加紧生产中\n")
			time.Sleep(time.Second)
		}
	}
END:
	return
}