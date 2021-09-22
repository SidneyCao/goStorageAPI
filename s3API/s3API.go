package main

import (
	"flag"
	"fmt"
)

//flag
var (
	bucket = flag.String("b", "", "bucket名 (默认为空)")
	files  = flag.String("f", "", "文件列表 (默认为空)")
	cache  = flag.String("c", "true", "是否缓存")
	prefix = flag.String("p", "", "需要移除的文件前缀 (默认为空)")
	thread = flag.Int("t", 5, "最大协程数")
)

func main() {
	fmt.Println("starting...")
}
