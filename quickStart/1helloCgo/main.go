package main

import "C"

// 通过 import "C" 语句启用 CGO 特性
func main() {
    // 通过 Go 内置的 println() 函数输出字符串，
    // 没有任何和 CGO 相关的代码。

    println("hello cgo")
}

// 虽然没有调用 CGO 相关的函数，但是 go build 命令在编译、链接阶段会启动 gcc 编译器，
// 这已经时一个完成的 CGO 程序了。
