package main

/*
#include <stdio.h>

static void SayHello(const char* s) {
    puts(s);
}
*/
import "C"

func main() {
    C.SayHello(C.CString("Hello, World\n"))
}

// 也可以将 SayHello() 函数放到当前目录下的一个 C 语言源文件中（扩展名必须是 .C）。
// 因为是编写在独立的 C 文件中，为了允许外部引用，所以需要去掉函数的 static 修饰符。
