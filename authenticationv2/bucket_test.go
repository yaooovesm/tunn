package authenticationv2

import (
	"fmt"
	"testing"
	"time"
)

func TestBucket(t *testing.T) {
	bucket := NewBucket(time.Second * 30)
	go func() {
		fmt.Println("协程1开启...")
		err := bucket.Occupy()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("协程1占用20秒...")
		time.Sleep(time.Second * 20)
		fmt.Println("协程1退出...")
		bucket.Leave()
	}()
	time.Sleep(time.Second)
	fmt.Println("主进程开启...")
	err := bucket.Occupy()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("主进程占用5秒")
	time.Sleep(time.Second * 5)
	fmt.Println("主进程退出...")
}
