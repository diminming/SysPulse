package common

import (
	"fmt"
	"runtime"
)

func PrintStackTrace() {
	// 创建一个缓冲区
	buf := make([]byte, 1<<16)
	// 将堆栈信息写入缓冲区，并打印到标准错误
	stackSize := runtime.Stack(buf, false)
	fmt.Printf("=== Recovered Stack Trace ===\n%s\n", buf[:stackSize])
}
