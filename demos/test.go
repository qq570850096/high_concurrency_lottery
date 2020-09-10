package main

//func main (){
//	syncmap := new(sync.Map)
//	ok := make(chan struct{})
//		go func() {				//开一个协程写map
//		for i := 0; i < 100; i++ {
//		syncmap.Store(i,i)
//	}
//	}()
//
//		go func() {				//开一个协程读map
//		for i := 0; i < 100; i++ {
//		fmt.Println(syncmap.Load(i))
//	}
//		ok <- struct{}{}
//	}()
//
//		//time.Sleep(time.Second * 20)
//		<-ok
//		// 猜猜会输出什么
//	}