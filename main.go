//package main
//
//import (
//	"fmt"
//	"math/rand"
//	"time"
//)
//
//type producer struct {
//	id    int
//	level int
//}
//
//type consumer struct {
//	id int
//}
//
//// 定义生产的方法
//func (p producer) produce(inputQ chan int, waitQ chan bool) {
//	var x int
//	x = rand.Intn(9) * p.level
//	for i := 0; i < 30; i++ { // 每个生产者生产30个产品
//		inputQ <- x // 将产品放入channel
//		fmt.Println("[===>] Producer ", p.id, " produced ", x)
//		sleepTime := rand.Intn(1000) + 1000
//		time.Sleep(time.Millisecond * time.Duration(sleepTime)) // 生产者每1~2秒生产一个产品
//	}
//	fmt.Println("[XXXX] Producer ", p.id, " quit.")
//	waitQ <- true // 生产者结束后在等待channel中放入一个信号
//}
//
//// 定义消费的方法
//func (c consumer) consume(inputQ chan int, waitQ chan bool) {
//	for {
//		select {
//		case x, ok := <-inputQ:
//			if ok { // 判断channel是否关闭了
//				fmt.Println("[<===] Consumer ", c.id, " consumed ", x)
//			}
//		default:
//			fmt.Println("[XXXX] Consumer ", c.id, " quit.")
//			waitQ <- true // 消费者结束后在等待channel中放入一个信号
//			return
//		}
//		sleepTime := rand.Intn(1000) + 2000
//		time.Sleep(time.Millisecond * time.Duration(sleepTime)) // 消费者每2~3秒消费一个产品
//	}
//}
//
//func main() {
//	var pLevels = [8]int{1, 10, 100, 100, 10000, 100000, 1000000, 10000000} // 定义不同的Level来区分生产者
//	const producerCount int = 5                                             // 生产者数量
//	const consumerCount int = 10                                            // 消费者数量
//	var inputQ = make(chan int, producerCount*4)                            // 生产队列，存放产品
//	var waitQ = make(chan bool)                                             // 等待队列，用于main方法等待所有消费者结束
//	var producers [producerCount]producer                                   // 生产者数组
//	var consumers [consumerCount]consumer                                   // 消费者数组
//	for i := 0; i < producerCount; i++ {
//		producers[i] = producer{i, pLevels[i]} // 初始化生产者id和Level
//		go producers[i].produce(inputQ, waitQ) // 开始生产
//	}
//
//	time.Sleep(time.Second * 5) // 积累一定产品后再启动消费者
//
//	for i := 0; i < consumerCount; i++ {
//		consumers[i] = consumer{i}             // 初始化消费者
//		go consumers[i].consume(inputQ, waitQ) // 开始消费
//	}
//
//	for i := 0; i < (producerCount + consumerCount); i++ {
//		<-waitQ // 等待所有生产者和消费者结束退出
//	}
//
//	fmt.Println("closed")
//	close(inputQ)
//	close(waitQ)
//}
