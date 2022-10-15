package main

/*
#include "hello.h"
*/
import "C"

func main() {
    C.SayHello(C.CString("Hello, module\n"))
}

// 模块化编程的核心是：面向接口编程（API）
