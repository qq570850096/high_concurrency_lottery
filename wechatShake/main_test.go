package wechatShake

import (
	"net/http"
	"sync"
	"testing"
)

func TestNewApp(t *testing.T) {
	main()
}

func BenchmarkNewApp(b *testing.B) {
	gof := func(wg *sync.WaitGroup){
		defer wg.Done()
		var resp *http.Response
		var err error
		for i:=0; i<2000; i++{
			//模拟一个get提交请求
			resp, err = http.Get("http://127.0.0.1:8080/lucky")
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close() //关闭连接
		}
	}
	var wg sync.WaitGroup

	for i:=0 ;i < 5 ; i++ {
		wg.Add(1)
		// 5个协程，应该是10000个数据
		go gof(&wg)
	}
	wg.Wait()
	// 实测执行后数据确实是一万个
}
