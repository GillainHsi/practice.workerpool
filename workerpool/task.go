package workerpool

import (
	"fmt"
	"time"
)

type Task struct {
	Id int
}

func (ref Task) DoSomething() {

	// 模擬task 執行時間
	<-time.After(time.Second * 1)

	fmt.Printf("%d task done\n", ref.Id)
}
