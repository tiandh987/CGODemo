package main

//static const char* cs = "hello";
import "C"

import "github.com/tiandh987/CGODemo/basics/cgoHelper/cgo_helper"

// 这段代码是不能正常工作的
func main() {
	// 当前 main 包引入的 C.cs 变量的类型是：当前 main 包的 CGO 构造的虚拟的 C 包下的 *char 类型（*main.C.char）;
	// 和 cgo_helper 包引入的 *C.char 类型（*cgo_helper.C.char）是不同的。
	cgo_helper.PrintCString(C.cs)
}

// CGO 将 “当前包” 引用的 C 语言符号都放到了虚拟的 C 包中，
// 同时当前包依赖的其它 Go 语言包内部可能也通过 CGO 引入了相似的虚拟 C 包，
// 但是，不同的 Go 语言包引入的虚拟的 C 包之间的类型是 “不能通用” 的。
