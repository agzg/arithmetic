package main

import "fmt"

func main() {
	lr := int(^uint(0)>>1) + 1
	fmt.Println(lr)
}
