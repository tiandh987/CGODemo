package main

//#include <stdio.h>
import "C"

func main() {
    C.puts(C.CString("Hello, World\n"))
}

// 通过 import "C" 语句启用 CGO 特性；
// 包含 C 语言的 <stdio.h> 头文件；
// 通过 cgo 包的 C.CString() 函数将 Go 语言字符串转换为 C 语言字符串；
// 通过 cgo 包的 C.puts() 函数向标准输出窗口打印转换后的字符串。

// 没有释放使用 C.CString 创建的 C 语言字符串会导致内存泄露。
// 但是对于这个小程序来说，没有问题，因为程序退出后操作系统会自动回收程序的所有资源。
