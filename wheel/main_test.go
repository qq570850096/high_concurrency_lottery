package wheel

import (
	"net/http"
	"sync"
	"testing"
)

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