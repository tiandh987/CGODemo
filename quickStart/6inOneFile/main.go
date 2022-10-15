package main

//void SayHello(char* s);
//void SayHelloV2(_GoString_ s);
import "C"

import "fmt"

func main() {
    C.SayHello(C.CString("Hello, in one file\n"))
    C.SayHelloV2("Hello, in one file V2\n")
}

//export SayHello
func SayHello(s *C.char) {
    fmt.Print(C.GoString(s))
}

//export SayHelloV2
func SayHelloV2(s string) {
    // 进一步以 Go 语言的思维来提炼 CGO 代码。
    // 分析发现 SayHello() 函数的参数如果可以直接使用 Go 字符串是最直接的。
    // 在 Go 1.10 中 CGO 新增加了一个 _GoString_ 预定义的 C 语言类型，用来表示 Go 语言字符串。
    fmt.Print(s)
}

// 执行过程：
//  Go main 函数 -> CGO 生成的 C 语言版本 SayHello() 桥接函数 -> Go 语言 SayHello() 函数
