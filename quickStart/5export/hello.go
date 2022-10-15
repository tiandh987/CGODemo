package main

// CGO 不仅用于 Go 语言中调用 C 语言函数，还可以用于导出 Go 语言函数给 C 语言函数调用。

import "C"

import "fmt"

//export SayHello
func SayHello(s *C.char) {
    fmt.Print(C.GoString(s))
}

// 通过 CGO 的 //export SayHello 指令将 Go 语言实现的函数 SayHello() 导出为 C 语言函数。

// 这里其实有两个版本的 SayHello：一个是 Go 语言环境的；另一个是 C 语言环境的。
// CGO 生成的 C 语言版本的 SayHello() 函数，最终会通过 “桥接代码” 调用 Go 语言版本的 SayHello() 函数。
