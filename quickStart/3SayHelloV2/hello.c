#include <stdio.h>

void SayHello(const char* s) {
    puts(s);
}

// 既然 SayHello() 函数已经放到独立的 C 文件中了，我们自然可以将对应的 C 文件编译打包为 静态库 或 动态库 文件供使用。

// 如果以 静态库 或 动态库 方式引用 SayHello() 函数，需要将对应的 C 源文件 “移出” 当前目录，
// CGO 构建程序会自动构建当前目录下的 C 源文件，从而导致 C 函数名冲突。