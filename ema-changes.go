package main

import (
	"fmt"
	"github.com/gogf/gf/text/gstr"
)

func main() {
	fmt.Println("Hello world")
	s := gstr.CaseSnake("Hello world")
	fmt.Println(s)
}
