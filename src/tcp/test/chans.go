package main

import (
	"fmt"
	"reflect"
	"sync"
	"time"
)

func chanTest() {
	chan1 := make(chan int, 1024)
	chan2 := make(chan int, 1024)
	for {
		select {
		case s1 := <-chan1:
			fmt.Println("chan1", s1)
		case s2 := <-chan2:
			fmt.Println("chan2", s2)
		default:
			chan1 <- 20
			chan2 <- 20
			time.Sleep(time.Millisecond * 100)
			break
		}

	}

}

// 互斥锁
// var x int64
// var wg sync.WaitGroup
// var lock sync.Mutex

// func add() {
// 	for i := 0; i < 5000; i++ {
// 		lock.Lock() // 加锁
// 		x = x + 1
// 		lock.Unlock() // 解锁
// 	}
// 	wg.Done()
// }
// func main() {
// 	wg.Add(2)
// 	go add()
// 	go add()
// 	wg.Wait()
// 	fmt.Println(x)
// }

// 读写互斥锁
// 互斥锁是完全互斥的，但是有很多实际的场景下是读多写少的，当我们并发的去读取一个资源不涉及资源修改的时候是没有必要加锁的，这种场景下使用读写锁是更好的一种选择。读写锁在Go语言中使用sync包中的RWMutex类型。

// 读写锁分为两种：读锁和写锁。当一个goroutine获取读锁之后，其他的goroutine如果是获取读锁会继续获得锁，如果是获取写锁就会等待；当一个goroutine获取写锁之后，其他的goroutine无论是获取读锁还是写锁都会等待。
// var (
// 	x      int64
// 	wg     sync.WaitGroup
// 	lock   sync.Mutex
// 	rwlock sync.RWMutex
// )

// func write() {
// 	// lock.Lock()   // 加互斥锁
// 	rwlock.Lock() // 加写锁
// 	x = x + 1
// 	time.Sleep(10 * time.Millisecond) // 假设读操作耗时10毫秒
// 	rwlock.Unlock()                   // 解写锁
// 	fmt.Println("wrrite", x)
// 	// lock.Unlock()                     // 解互斥锁
// 	wg.Done()
// }

// func read() {
// 	// lock.Lock()                  // 加互斥锁
// 	rwlock.RLock()               // 加读锁
// 	time.Sleep(time.Millisecond) // 假设读操作耗时1毫秒
// 	fmt.Println("read", x)
// 	rwlock.RUnlock() // 解读锁
// 	// lock.Unlock()                // 解互斥锁
// 	wg.Done()
// }

// func main() {
// 	start := time.Now()
// 	for i := 0; i < 10; i++ {
// 		wg.Add(1)
// 		go write()
// 	}

// 	for i := 0; i < 1000; i++ {
// 		wg.Add(1)
// 		go read()
// 	}

// 	wg.Wait()
// 	end := time.Now()
// 	fmt.Println(end.Sub(start))
// }

// 定义结构体
type User struct {
	Id   int
	Name string
	Age  int
}

func (u User) Hello(name string) {
	fmt.Println("Hello：", name)
}
func (u *User) TestHello(name string) {
	fmt.Println("TestHello：", name)
}

func ReflectTest() {
	if true {
		u := User{1, "5lmh.com", 20}
		v := reflect.ValueOf(u)
		m := v.MethodByName("Hello")
		args := []reflect.Value{reflect.ValueOf("6666")}
		m.Call(args)
	}
	if true {
		u := &User{1, "5lmh.com", 20}

		// 获取User实例的指针的反射值
		v := reflect.ValueOf(u)

		// 因为TestHello是一个指针接收器的方法，我们需要获取指针指向的值的方法
		m := v.MethodByName("TestHello")

		// 准备参数
		args := []reflect.Value{reflect.ValueOf("6666")}

		// 调用方法
		m.Call(args)
	}
}

var wg sync.WaitGroup
var mutes sync.Mutex
var c int = 0

func LockTest() {
	mutes.Lock()
	defer mutes.Unlock()
	c = c + 1
	wg.Done()
}

func main() {
	for i := 0; i < 1000000; i++ {
		wg.Add(1)
		go LockTest()
	}
	wg.Wait()
	fmt.Println("c", c)
}
