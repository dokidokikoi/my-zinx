package ztimer

import (
	"fmt"
	"testing"
	"time"
)

// 定义一个超时函数
func myFunc(v ...interface{}) {
	fmt.Printf("time %d No.%d function calld. delay %d second(s)\n", time.Now().Second(), v[0].(int), v[1].(int))
}

func TestTimer(t *testing.T) {

	for i := 0; i < 5; i++ {
		//go func(i int) {
		NewTimerAfter(NewDelayFunc(myFunc, []interface{}{i, 2 * i}), time.Duration(2*i)*time.Second).Run()
		//}(i)
	}

	//主进程等待其他go，由于Run()方法是用一个新的go承载延迟方法，这里不能用waitGroup
	time.Sleep(1 * time.Minute)
}
