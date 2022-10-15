package main

//void SayHello(const char* s);
import "C"

func main() {
    C.SayHello(C.CString("Hello, World V2\n"))
}

// 使用 go build 命令编译
