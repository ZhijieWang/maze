package main

import (
	"fmt"
	"maze/common"
)

func main() {
	fmt.Println("hello")
	common.CreateWorld(2, common.NewBasicTaskManager())
}
