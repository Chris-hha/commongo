package main

// 下面的示例代码演示了如何使用 errgroup 包来处理多个子任务 goroutine 中可能返回的 error。
import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"github.com/pkg/errors"
)


func ErrorGroupDemo() error {
	g := new(errgroup.Group) // 创建等待组（类似sync.WaitGroup）
	var testSlice = []interface{}{
		123,
		true,
		"http://pkg.go.dev",
		"http://www.test.com",
		"http://www.liwenzhou.com",
	}
	for _, item := range testSlice {
		// 启动一个goroutine去获取url内容
		g.Go(func(item interface{}) func() error {
			return func() error {
				_, ok := item.(string)
				if ok {
					fmt.Printf("%v is string type\n", item)
					return nil
				}
				return errors.Errorf("%v is not string", item)
			}
		}(item))
	}
	if err := g.Wait(); err != nil {
		// 处理可能出现的错误
		panic(err)
	}
	fmt.Println("所有goroutine均成功")
	return nil
}

func main(){
	err := ErrorGroupDemo()
	if err == nil {
		fmt.Println("exec success")
	}
}
