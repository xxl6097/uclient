package main

import (
	"fmt"
	"github.com/xxl6097/uclient/internal/u"
	"os"
)

func main() {
	args := os.Args[1:]
	fmt.Println(os.Args)
	fmt.Println(args)
	a := u.Ping(args[0])
	fmt.Println(a)

}
