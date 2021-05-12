// 原子性测试脚本
// ++ 的操作是非原子操作，在执行过程中会由于多个goroutine并发操作，使结果跟预期值不一致
// 利用atomic原语，可以实现原子操作，使执行结果符合预期
// Reference：https://golang.org/pkg/sync/atomic/

package main

import (
	"flag"
	"fmt"
	"sync"
	"sync/atomic"
)

var wg sync.WaitGroup
var count int //计数器

func main() {
	//定一个变量，代表程序是否以原子操作方式运行
	var isAtomic bool
	//加入命令行参数 -a
	flag.BoolVar(&isAtomic, "a", false, "是否想使用[-a]来指定程序以原子操作运行")
	flag.Parse()

	fmt.Printf("isAtomic:%v \n", isAtomic)

	if isAtomic {
		showCountAtomic()
	} else {
		showCount()
	}
}

//非原子统计
func showCount() {
	for i := 1; i <= 1000000; i++ {
		wg.Add(1)
		go incr(&count)
	}
	wg.Wait()
	fmt.Println("-------")
	fmt.Println(count)
}

//非原子++操作
func incr(k *int) {
	defer wg.Done()
	*k ++
}

//原子统计
func showCountAtomic() {
	var kp int32
	kp = (int32)(count)
	for i := 1; i <= 100000; i++ {
		wg.Add(1)
		go incrAtomic(&kp)
	}
	wg.Wait()
	count = int(kp)
	fmt.Println("-------")
	fmt.Println(count)
}

//原子++操作
func incrAtomic(k *int32) {
	defer wg.Done()
	atomic.AddInt32(k, 1)
}
